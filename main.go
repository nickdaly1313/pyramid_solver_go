package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"pyramid_solver_go_local/game"   // <--- Ensure this path is correct
	"pyramid_solver_go_local/solver" // <--- Ensure this path is correct
	
)

func main() {
	fmt.Println("Welcome to the Pyramid Stone Puzzle Solver!")
	reader := bufio.NewReader(os.Stdin)

	for { // Main loop to solve multiple puzzles
		gameInstance := game.NewPuzzleGame()

		fmt.Println("\nHow would you like to set up the puzzle?")
		fmt.Println("1. Enter pyramid and draw pile values manually")
		fmt.Println("2. Use the example puzzle from our discussion")
		fmt.Print("Enter your choice (1 or 2): ")

		choiceStr, _ := reader.ReadString('\n')
		choiceStr = strings.TrimSpace(choiceStr)
		choice, _ := strconv.Atoi(choiceStr)

		var pyramidStones []int
		var drawPileStones []int
		var err error

		if choice == 1 {
			pyramidStones, err = getPyramidInput(reader)
			if err != nil {
				fmt.Printf("Error getting pyramid input: %v\n", err)
				// Optionally, ask if they want to try again or exit
				fmt.Print("Try again? (y/n): ")
				retry, _ := reader.ReadString('\n')
				if strings.TrimSpace(strings.ToLower(retry)) == "y" {
					continue // Restart the loop
				}
				break // Exit if not retrying
			}
			drawPileStones, err = getDrawPileInput(reader) // Updated call
			if err != nil {
				fmt.Printf("Error getting draw pile input: %v\n", err)
				// Optionally, ask if they want to try again or exit
				fmt.Print("Try again? (y/n): ")
				retry, _ := reader.ReadString('\n')
				if strings.TrimSpace(strings.ToLower(retry)) == "y" {
					continue // Restart the loop
				}
				break // Exit if not retrying
			}
			err = gameInstance.SetupCustomGame(pyramidStones, drawPileStones)
		} else {
			fmt.Println("\nUsing the example puzzle from our discussion...")
			pyramidStones = []int{12, 10, 11, 6, 11, 7, 12, 11, 5, 1, 4, 1, 4, 5, 10, 8, 11, 9, 7, 2, 9, 6, 2, 13, 9, 10, 12, 13}
			drawPileStones = []int{6, 3, 8, 9, 3, 10, 2, 13, 6, 7, 1, 13, 12, 4, 1, 2, 3, 8, 5, 3, 5, 7, 3, 8}
			err = gameInstance.SetupCustomGame(pyramidStones, drawPileStones)
		}

		if err != nil {
			fmt.Printf("Error setting up game: %v\n", err)
			fmt.Print("Try again? (y/n): ")
			retry, _ := reader.ReadString('\n')
			if strings.TrimSpace(strings.ToLower(retry)) == "y" {
				continue // Restart the loop
			}
			break // Exit if not retrying
		}

		fmt.Println("\nInitial Game State:")
		gameInstance.PrintState()

		puzzleSolver := solver.NewPuzzleSolver(gameInstance)

        fmt.Println("\nSolving puzzle...")
        bestMoves, bestScore := puzzleSolver.SolveMonteCarlo(100000) // Set iterations to 100,000


		fmt.Printf("\nBest solution found - Score: %d, Moves: %d\n", bestScore, len(bestMoves))

		solutionText := formatSolution(bestMoves, bestScore)
		fmt.Println("\n" + solutionText)

		// Ask to solve another puzzle
		fmt.Print("\nSolve another pyramid? (y/n): ")
		anotherPuzzleStr, _ := reader.ReadString('\n')
		anotherPuzzle := strings.TrimSpace(strings.ToLower(anotherPuzzleStr))
		if anotherPuzzle != "n" {
			break // Exit the loop if the user doesn't want to continue
		}
	}
}

func getPyramidInput(reader *bufio.Reader) ([]int, error) {
    fmt.Println("\n=== PYRAMID INPUT ===")
    fmt.Println("Enter 28 characters (a-u) for the pyramid stones.")
    fmt.Println("a=1, s=2, d=3, f=4, g=5, h=6, j=7, k=8, l=9, r=10, t=11, y=12, u=13")
    fmt.Println("No spaces between characters.")

    // ... (Print pyramid structure - this can remain the same)

    for {
        fmt.Print("\nEnter 28 characters: ")
        inputStr, _ := reader.ReadString('\n')
        inputStr = strings.TrimSpace(inputStr)

        if len(inputStr) != game.TotalPyramidStones {
            fmt.Println("Invalid input length. Please enter 28 characters.")
            continue
        }

        pyramidStones := make([]int, game.TotalPyramidStones)
        for i, char := range inputStr {
            stone, err := charToInt(char)
            if err != nil {
                fmt.Println(err)
                continue // Go to the next input attempt
            }
            pyramidStones[i] = stone
        }
        return pyramidStones, nil
    }
}


func getDrawPileInput(reader *bufio.Reader) ([]int, error) {
    fmt.Println("\n=== DRAW PILE INPUT ===")
    fmt.Println("Enter 24 characters (a-u) for the draw pile stones.")
    fmt.Println("a=1, s=2, d=3, f=4, g=5, h=6, j=7, k=8, l=9, r=10, t=11, y=12, u=13")
    fmt.Println("No spaces between characters.")

    for {
        fmt.Print("\nEnter 24 characters: ")
        inputStr, _ := reader.ReadString('\n')
        inputStr = strings.TrimSpace(inputStr)

        if len(inputStr) != game.MaxDrawPileSegments*game.StonesPerSegment {
            fmt.Println("Invalid input length. Please enter 24 characters.")
            continue
        }

        drawPileStones := make([]int, game.MaxDrawPileSegments*game.StonesPerSegment)
        for i, char := range inputStr {
            stone, err := charToInt(char)
            if err != nil {
                fmt.Println(err)
                continue // Go to the next input attempt
            }
            drawPileStones[i] = stone
        }
        return drawPileStones, nil
    }
}


// charToInt converts a character to its corresponding integer value.
func charToInt(char rune) (int, error) {
    switch char {
    case 'a': return 1, nil
    case 's': return 2, nil
    case 'd': return 3, nil
    case 'f': return 4, nil
    case 'g': return 5, nil
    case 'h': return 6, nil
    case 'j': return 7, nil
    case 'k': return 8, nil
    case 'l': return 9, nil
    case 'r': return 10, nil
    case 't': return 11, nil
    case 'y': return 12, nil
    case 'u': return 13, nil
    default: return -1, fmt.Errorf("invalid character: %c", char)
    }
}

// formatSolution formats the solution into a readable string.
func formatSolution(moves []game.Move, score int) string {
    var sb strings.Builder
    sb.WriteString(fmt.Sprintf("Final Score: %d\n", score))
    sb.WriteString("\nStep-by-Step Solution:\n")

    for _, move := range moves { // Removed the index 'i' from the loop
        // sb.WriteString(fmt.Sprintf("%d. ", i+1))  // Remove this line - no longer printing the number

        switch {
        case move.Source == "DRAW" && move.Destination == "DRAW":
            sb.WriteString("DRAW\n")
        case move.Destination == "HOLD":
            sb.WriteString(fmt.Sprintf("Move %s to HOLD\n", move.Source))
        case move.Source == "HOLD":
            sb.WriteString(fmt.Sprintf("Match HOLD and %s\n", move.Destination))
        case move.Destination == "SMASH":
            sb.WriteString(fmt.Sprintf("Smash %s\n", move.Source))
        case move.Source == "DRW1":
            if move.Destination == "SMASH" {
                sb.WriteString("Smash DRW1\n")
            } else {
                sb.WriteString(fmt.Sprintf("Move DRW1 to %s\n", move.Destination))
            }
        case move.Destination == "DRW1":
            sb.WriteString(fmt.Sprintf("Match %s to DRW1\n", move.Source))
        default:
            sb.WriteString(fmt.Sprintf("Match %s and %s\n", move.Source, move.Destination))
        }
    }
    return sb.String()
}
