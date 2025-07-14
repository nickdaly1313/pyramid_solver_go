package solver

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"pyramid_solver_go_local/game"
	"pyramid_solver_go_local/utils"
)

// PuzzleSolver manages the Monte Carlo simulation.
type PuzzleSolver struct {
	originalGame *game.PuzzleGame
	bestScore    int
	bestMoves    []game.Move
}

// NewPuzzleSolver creates a new PuzzleSolver.
func NewPuzzleSolver(originalGame *game.PuzzleGame) *PuzzleSolver {
	return &PuzzleSolver{
		originalGame: originalGame,
		bestScore:    -1,
		bestMoves:    []game.Move{},
	}
}

// --- Structs for Parallel Processing ---
type Job struct {
	NumSimulations int
	Seed           int64
}
type Result struct {
	Score int
	Moves []game.Move
}

// --- Worker Function (Updated for "Double Reset" Pattern) ---
func (s *PuzzleSolver) worker(jobs <-chan Job, results chan<- Result) {
	r := rand.New(rand.NewSource(0))

	// Each worker allocates TWO game objects and reuses them.
	simulatedGame := s.originalGame.DeepCopy() // For the main simulation
	tempGame := s.originalGame.DeepCopy()      // A reusable "scratchpad" for testing moves

	for job := range jobs {
		r.Seed(job.Seed)

		localBestScore := -1
		var localBestMoves []game.Move

		for i := 0; i < job.NumSimulations; i++ {
			simulatedGame.Reset(s.originalGame)
			movesMade := []game.Move{}

			for !simulatedGame.IsSolved() {
				possibleMoves := s.getPossibleMovesForSimulation(simulatedGame)
				if len(possibleMoves) == 0 {
					break
				}

				var chosenMove game.Move
				if simulatedGame.CalculateScore() == 0 {
					chosenMove = possibleMoves[r.Intn(len(possibleMoves))]
				} else {
					matchingMoves := []game.Move{}
					for _, move := range possibleMoves {
						// *** THE SECOND KEY PERFORMANCE FIX IS HERE ***
						// Reset the tempGame to the current simulation state.
						tempGame.Reset(simulatedGame)
						if tempGame.MakeMove(move.Source, move.Destination) {
							matchingMoves = append(matchingMoves, move)
						}
					}
					if len(matchingMoves) > 0 && r.Float64() < 0.8 {
						chosenMove = matchingMoves[r.Intn(len(matchingMoves))]
					} else {
						chosenMove = possibleMoves[r.Intn(len(possibleMoves))]
					}
				}
				simulatedGame.MakeMove(chosenMove.Source, chosenMove.Destination)
				movesMade = append(movesMade, chosenMove)
				if len(movesMade) > 200 {
					break
				}
			}

			finalScore := simulatedGame.CalculateScore()
			if finalScore > localBestScore {
				localBestScore = finalScore
				localBestMoves = make([]game.Move, len(movesMade))
				copy(localBestMoves, movesMade)
			}
		}
		results <- Result{Score: localBestScore, Moves: localBestMoves}
	}
}

// --- Manager Function (Updated for Batching) ---
func (s *PuzzleSolver) SolveMonteCarlo(iterations int) ([]game.Move, int) {
	fmt.Printf("Running %d simulations in parallel...\n", iterations)

	numWorkers := runtime.NumCPU()
	fmt.Printf("Utilizing %d CPU cores as workers.\n", numWorkers)

	jobs := make(chan Job, numWorkers)
	results := make(chan Result, numWorkers)

	for w := 0; w < numWorkers; w++ {
		go s.worker(jobs, results)
	}

	simsPerWorker := iterations / numWorkers
	baseSeed := time.Now().UnixNano()
	for w := 0; w < numWorkers; w++ {
		numSims := simsPerWorker
		if w == numWorkers-1 {
			numSims += iterations % numWorkers
		}
		if numSims > 0 {
			jobs <- Job{NumSimulations: numSims, Seed: baseSeed + int64(w)}
		}
	}
	close(jobs)

	fmt.Println("All jobs distributed. Collecting results...")
	activeWorkers := numWorkers
	// Adjust loop to handle cases where numSims was 0 for some workers
	for w := 0; w < numWorkers; w++ {
		if (iterations/numWorkers > 0) || (w < iterations%numWorkers) {
			result := <-results
			activeWorkers--
			fmt.Printf("\rResult received. Waiting for %d more workers...", activeWorkers)
			if result.Score > s.bestScore {
				s.bestScore = result.Score
				s.bestMoves = result.Moves
			}
		}
	}
	fmt.Println("\nCollection complete.")

	return s.bestMoves, s.bestScore
}

// getPossibleMovesForSimulation remains unchanged.
func (s *PuzzleSolver) getPossibleMovesForSimulation(g *game.PuzzleGame) []game.Move {
	moves := []game.Move{{Source: "DRAW", Destination: "DRAW"}}
	accessiblePositions := g.GetAccessiblePositions()
	for _, pos := range accessiblePositions {
		row, col, _ := utils.StringToIndices(pos)
		if g.PyramidValue(row, col) == 13 {
			moves = append(moves, game.Move{Source: pos, Destination: "SMASH"})
		}
	}
	for i, pos1 := range accessiblePositions {
		row1, col1, _ := utils.StringToIndices(pos1)
		stone1 := g.PyramidValue(row1, col1)
		if stone1 == 13 {
			continue
		}
		for j := i + 1; j < len(accessiblePositions); j++ {
			pos2 := accessiblePositions[j]
			row2, col2, _ := utils.StringToIndices(pos2)
			stone2 := g.PyramidValue(row2, col2)
			if g.IsMatchingPair(stone1, stone2) {
				moves = append(moves, game.Move{Source: pos1, Destination: pos2})
			}
		}
		if g.HoldValue() != -1 && g.IsMatchingPair(stone1, g.HoldValue()) {
			moves = append(moves, game.Move{Source: pos1, Destination: "HOLD"})
		}
		drw1Stone := g.GetCurrentDrawStone()
		if drw1Stone != -1 && g.IsMatchingPair(stone1, drw1Stone) {
			moves = append(moves, game.Move{Source: pos1, Destination: "DRW1"})
		}
		if g.HoldValue() == -1 && stone1 != 13 {
			moves = append(moves, game.Move{Source: pos1, Destination: "HOLD"})
		}
	}
	if g.HoldValue() != -1 && g.GetCurrentDrawStone() != -1 && g.IsMatchingPair(g.HoldValue(), g.GetCurrentDrawStone()) {
		moves = append(moves, game.Move{Source: "HOLD", Destination: "DRW1"})
	}
	if g.HoldValue() == -1 && g.GetCurrentDrawStone() != -1 && g.GetCurrentDrawStone() != 13 {
		moves = append(moves, game.Move{Source: "DRW1", Destination: "HOLD"})
	}
	if g.GetCurrentDrawStone() == 13 {
		moves = append(moves, game.Move{Source: "DRW1", Destination: "SMASH"})
	}
	return moves
}

// --- Helper functions for accessing game state ---
func (s *PuzzleSolver) IsMatchingPair(stone1, stone2 int) bool {
	return s.originalGame.IsMatchingPair(stone1, stone2)
}
func (s *PuzzleSolver) PyramidValue(row, col int) int {
	return s.originalGame.PyramidValue(row, col)
}
func (s *PuzzleSolver) HoldValue() int {
	return s.originalGame.HoldValue()
}
func (s *PuzzleSolver) GetCurrentDrawStone() int {
	return s.originalGame.GetCurrentDrawStone()
}
func (s *PuzzleSolver) Redraws() int {
	return s.originalGame.Redraws()
}
func (s *PuzzleSolver) NumActiveSegments() int {
	return s.originalGame.NumActiveSegments()
}
