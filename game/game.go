package game


import (
   "fmt"
   "math"
   "strings"


   "pyramid_solver_go_local/utils" // <--- Ensure this path is correct for your repo
)


const (
   MaxDrawPileSegments = 8
   StonesPerSegment    = 3
   TotalPyramidStones  = 28
   MaxPyramidRows      = 7
   MaxPyramidCols      = 7 // Max columns in any row (for array sizing)
)


// Move represents a game move.
type Move struct {
   Source      string
   Destination string
}


// PuzzleGame represents the state of the pyramid puzzle game.
type PuzzleGame struct {
   pyramid             [MaxPyramidRows][MaxPyramidCols]int // Stores stone values, -1 for empty
   hold                int                                 // -1 for empty
   drawPile            [MaxDrawPileSegments][]int          // Array of slices for segments
   currentSegment      int
   moves               []Move
   matches             int
   streak              int
   streakBonus         int
   redraws             int
   timeRemaining       int // Fixed at 120 for scoring
   numActiveSegments   int // Actual number of active segments in drawPile
}




// NewPuzzleGame creates and initializes a new PuzzleGame.
func NewPuzzleGame() *PuzzleGame {
   game := &PuzzleGame{
       hold:          -1, // -1 indicates empty hold
       timeRemaining: 120,
   }
   game.initializePyramid()
   return game
}


// initializePyramid sets up the initial empty pyramid structure.
func (g *PuzzleGame) initializePyramid() {
   for rowIdx := 0; rowIdx < MaxPyramidRows; rowIdx++ {
       for colIdx := 0; colIdx < MaxPyramidCols; colIdx++ { // Initialize all possible slots
           g.pyramid[rowIdx][colIdx] = -1 // -1 indicates empty
       }
   }
}


// SetupRandomGame sets up a random game configuration.
func (g *PuzzleGame) SetupRandomGame() {
   stones := []int{}
   for i := 1; i <= 12; i++ {
       for j := 0; j < 4; j++ {
           stones = append(stones, i)
       }
   }
   for i := 0; i < 4; i++ {
       stones = append(stones, 13)
   }
   stones = utils.ShuffleArray(stones)


   // Fill the pyramid
   pyramidPositions := g.getAllPyramidPositions()
   for i := 0; i < TotalPyramidStones; i++ {
       g.pyramid[pyramidPositions[i][0]][pyramidPositions[i][1]] = stones[i]
   }


   // Fill the draw pile
   // Corrected loop header: segmentIdx is declared once
   for segmentIdx := 0; segmentIdx < MaxDrawPileSegments; segmentIdx++ {
       start := TotalPyramidStones + segmentIdx*StonesPerSegment
       end := start + StonesPerSegment
       if end > len(stones) { // Handle cases where not enough stones for full segments
           end = len(stones)
       }
       g.drawPile[segmentIdx] = stones[start:end]
   }
   g._trimEmptySegments()
}


// SetupCustomGame sets up the game with user-provided values.
func (g *PuzzleGame) SetupCustomGame(pyramidStones, drawPileStones []int) error {
   if len(pyramidStones) != TotalPyramidStones {
       return fmt.Errorf("pyramid must have %d stones, got %d", TotalPyramidStones, len(pyramidStones))
   }
   if len(drawPileStones) > MaxDrawPileSegments*StonesPerSegment { // Allow fewer than 24 if user provides
       return fmt.Errorf("draw pile must have at most %d stones, got %d", MaxDrawPileSegments*StonesPerSegment, len(drawPileStones))
   }


   pyramidPositions := g.getAllPyramidPositions()
   for i := 0; i < TotalPyramidStones; i++ {
       g.pyramid[pyramidPositions[i][0]][pyramidPositions[i][1]] = pyramidStones[i]
   }


   for segmentIdx := 0; segmentIdx < MaxDrawPileSegments; segmentIdx++ {
       start := segmentIdx * StonesPerSegment
       end := start + StonesPerSegment
       if end > len(drawPileStones) {
           end = len(drawPileStones)
       }
       g.drawPile[segmentIdx] = drawPileStones[start:end]
   }
   g._trimEmptySegments()
   return nil
}


// getAllPyramidPositions returns a slice of [row, col] indices for all pyramid positions.
// Row 0 is 'A' (bottom), Row 6 is 'G' (top).
func (g *PuzzleGame) getAllPyramidPositions() [][2]int {
   positions := make([][2]int, 0, TotalPyramidStones)
   for rowIdx := 0; rowIdx < MaxPyramidRows; rowIdx++ { // Iterate 0 (A) to 6 (G)
       for colIdx := 0; colIdx < utils.PyramidRowSizes[rowIdx]; colIdx++ {
           positions = append(positions, [2]int{rowIdx, colIdx})
       }
   }
   return positions
}


// IsMatchingPair checks if two stones form a valid matching pair.
func (g *PuzzleGame) IsMatchingPair(stone1, stone2 int) bool {
   if stone1 == -1 || stone2 == -1 || stone1 == 13 || stone2 == 13 {
       return false
   }
   validPairs := [][2]int{{1, 2}, {3, 4}, {5, 6}, {7, 8}, {9, 10}, {11, 12}}
   for _, pair := range validPairs {
       if (stone1 == pair[0] && stone2 == pair[1]) || (stone1 == pair[1] && stone2 == pair[0]) {
           return true
       }
   }
   return false
}


// IsAccessible checks if a pyramid position is accessible.
// Row 0 is 'A' (bottom), Row 6 is 'G' (top).
func (g *PuzzleGame) IsAccessible(rowIdx, colIdx int) bool {
   if g.pyramid[rowIdx][colIdx] == -1 { // Position is empty
       return false
   }
   if rowIdx == 0 { // Row A (bottom, index 0) is always accessible
       return true
   }
   // Check if the two positions below are empty
   return g.pyramid[rowIdx-1][colIdx] == -1 && g.pyramid[rowIdx-1][colIdx+1] == -1
}


// GetAccessiblePositions returns a slice of accessible pyramid positions as strings.
func (g *PuzzleGame) GetAccessiblePositions() []string {
   accessible := []string{}
   for rowIdx := 0; rowIdx < MaxPyramidRows; rowIdx++ {
       for colIdx := 0; colIdx < utils.PyramidRowSizes[rowIdx]; colIdx++ {
           if g.IsAccessible(rowIdx, colIdx) {
               pos, _ := utils.IndicesToString(rowIdx, colIdx) // Error handling omitted for brevity
               accessible = append(accessible, pos)
           }
       }
   }
   return accessible
}


// game/game.go (GetCurrentDrawStone with backfilling and debug prints)


func (g *PuzzleGame) GetCurrentDrawStone() int {
  


   if len(g.drawPile[g.currentSegment]) > 0 {
       stone := g.drawPile[g.currentSegment][len(g.drawPile[g.currentSegment])-1]
       return stone
   }


   // Backfill only if the current segment is empty
   if len(g.drawPile[g.currentSegment]) == 0 {
       for seg := g.currentSegment - 1; seg >= 0; seg-- {
           if len(g.drawPile[seg]) > 0 {
               stone := g.drawPile[seg][len(g.drawPile[seg])-1]
               return stone
           }
       }
   }


   return -1 // No stone available
}




// MakeMove performs a move in the game. Returns true if a stone was cleared (match or smash).
func (g *PuzzleGame) MakeMove(source, destination string) bool {
  g.moves = append(g.moves, Move{Source: source, Destination: destination})




  if source == "DRAW" && destination == "DRAW" {
       g.currentSegment++
       if g.currentSegment >= g.numActiveSegments { // Check if we've reached the end of the draw pile
           g.currentSegment = 0                      // Reset to the first segment
           g.redraws++
           g._redistributeDrawPile()                // Redistribute and update numActiveSegments
           g._trimEmptySegments()
       }
       g.streak = 0
       return false
   }




  if destination == "SMASH" {
      var stone int
      if source == "DRW1" {
          stone = g.GetCurrentDrawStone()  // Get the stone value (could be from a previous segment)
          if stone == 13 {
           // Remove the smashed stone (handling backfilling)
               if len(g.drawPile[g.currentSegment]) > 0 {
                   g.drawPile[g.currentSegment] = g.drawPile[g.currentSegment][:len(g.drawPile[g.currentSegment])-1]
               } else { // Backfill case
                   for seg := g.currentSegment - 1; seg >= 0; seg-- {
                       if len(g.drawPile[seg]) > 0 {
                           g.drawPile[seg] = g.drawPile[seg][:len(g.drawPile[seg])-1]
                           g._trimEmptySegments() // Important: Update numActiveSegments
                           break // Stop after removing from one segment
                       }
                   }
               }
          } else {
              return false
          }
      } else {
          row, col, _ := utils.StringToIndices(source) // Error handling omitted for brevity
          stone = g.pyramid[row][col]
          if stone == 13 {
              g.pyramid[row][col] = -1
          } else {
              return false
          }
      }
      g.matches++
      g.streak++
      if g.streak > 1 { // Only add bonus if streak is 2 or more
          if g.streak >= 5 { // Cap at 200 for 5th match and greater
              g.streakBonus += 200
          } else { // For 2nd, 3rd, 4th match
              g.streakBonus += (g.streak - 1) * 50
          }
      }
      return true // Stone cleared
  }




  var sourceStone int
  if source == "HOLD" {
      sourceStone = g.hold
  } else if source == "DRW1" {
      sourceStone = g.GetCurrentDrawStone()
  } else {
      row, col, _ := utils.StringToIndices(source) // Error handling omitted for brevity
      sourceStone = g.pyramid[row][col]
  }




  // --- CRITICAL FIX: Determine if it's a match *first* ---
  // Get the potential destination stone for matching purposes
  var potentialMatchDestStone int
  if destination == "HOLD" { // If destination is HOLD, the 'other' stone for matching would be the current hold content
      potentialMatchDestStone = g.hold
  } else if destination == "DRW1" {
      potentialMatchDestStone = g.GetCurrentDrawStone()
  } else {
      row, col, _ := utils.StringToIndices(destination)
      potentialMatchDestStone = g.pyramid[row][col]
  }




  // If it's a matching pair, perform the match
  if g.IsMatchingPair(sourceStone, potentialMatchDestStone) {
      // Remove source stone
      if source == "DRW1" {
          if len(g.drawPile[g.currentSegment]) > 0 { // Remove from current segment if not empty
              g.drawPile[g.currentSegment] = g.drawPile[g.currentSegment][:len(g.drawPile[g.currentSegment])-1]
       } else { // Backfill case: Remove from previous non-empty segment
               for seg := g.currentSegment - 1; seg >= 0; seg-- {
                   if len(g.drawPile[seg]) > 0 {
                       g.drawPile[seg] = g.drawPile[seg][:len(g.drawPile[seg])-1]
                       g._trimEmptySegments() // Important: Update numActiveSegments if a segment becomes empty
                       break // Stop after removing from one segment
                   }
               }
           }
      } else if source == "HOLD" {
          g.hold = -1
      } else {
          row, col, _ := utils.StringToIndices(source) // Error handling omitted for brevity
          g.pyramid[row][col] = -1
      }




      // Remove destination stone
      if destination == "DRW1" {
          if len(g.drawPile[g.currentSegment]) > 0 { // Remove from current segment if not empty
              g.drawPile[g.currentSegment] = g.drawPile[g.currentSegment][:len(g.drawPile[g.currentSegment])-1]
           } else { // Backfill case: Remove from previous non-empty segment
               for seg := g.currentSegment - 1; seg >= 0; seg-- {
                   if len(g.drawPile[seg]) > 0 {
                       g.drawPile[seg] = g.drawPile[seg][:len(g.drawPile[seg])-1]
                       g._trimEmptySegments() // Important: Update numActiveSegments if a segment becomes empty
                       break // Stop after removing from one segment
                   }
               }
          }
      } else if destination == "HOLD" {
          g.hold = -1
      } else {
          row, col, _ := utils.StringToIndices(destination) // Error handling omitted for brevity
          g.pyramid[row][col] = -1
      }




      g.matches++
      g.streak++
      if g.streak > 1 { // Only add bonus if streak is 2 or more
          if g.streak >= 5 { // Cap at 200 for 5th match and greater
              g.streakBonus += 200
          } else { // For 2nd, 3rd, 4th match
              g.streakBonus += (g.streak - 1) * 50
          }
      }
      return true // Stone cleared
  }




  // --- If it's NOT a match, then check if it's a valid move to HOLD ---
  if destination == "HOLD" {
      if g.hold != -1 { // Hold is already occupied, so this is an invalid move
          return false
      }
      g.hold = sourceStone
      if source == "DRW1" {
          if len(g.drawPile[g.currentSegment]) > 0 { // Remove from current segment if not empty
              g.drawPile[g.currentSegment] = g.drawPile[g.currentSegment][:len(g.drawPile[g.currentSegment])-1]
           } else { // Backfill case: Remove from previous non-empty segment
               for seg := g.currentSegment - 1; seg >= 0; seg-- {
                   if len(g.drawPile[seg]) > 0 {
                       g.drawPile[seg] = g.drawPile[seg][:len(g.drawPile[seg])-1]
                       g._trimEmptySegments() // Important: Update numActiveSegments if a segment becomes empty
                       break // Stop after removing from one segment
                   }
               }
          }
      } else {
          row, col, _ := utils.StringToIndices(source) // Error handling omitted for brevity
          g.pyramid[row][col] = -1
      }
      return false // Not a stone-clearing move
  }




  // If it's neither a match nor a valid move to HOLD, it's a non-clearing move
  g.streak = 0 // Streak resets for any other non-clearing move
  return false // Not a stone-clearing move
}




// _redistributeDrawPile redistributes the remaining stones in the draw pile into new segments.
func (g *PuzzleGame) _redistributeDrawPile() {
   allStones := []int{}
   for _, segment := range g.drawPile {
       allStones = append(allStones, segment...)
   }


   g.drawPile = [MaxDrawPileSegments][]int{} // Clear draw pile
   g.numActiveSegments = 0


   segmentIdx := 0
   for _, stone := range allStones {
       if segmentIdx < MaxDrawPileSegments {
           if len(g.drawPile[segmentIdx]) < StonesPerSegment {
               g.drawPile[segmentIdx] = append(g.drawPile[segmentIdx], stone)
           } else {
               segmentIdx++
               if segmentIdx < MaxDrawPileSegments {
                   g.drawPile[segmentIdx] = append(g.drawPile[segmentIdx], stone)
               }
           }
       }
   }
   g._trimEmptySegments()
}


// _trimEmptySegments removes empty segments from the end of the draw pile.
func (g *PuzzleGame) _trimEmptySegments() {
   g.numActiveSegments = MaxDrawPileSegments
   for g.numActiveSegments > 0 && len(g.drawPile[g.numActiveSegments-1]) == 0 {
       g.numActiveSegments--
   }
}


// IsSolved checks if the puzzle is solved (pyramid is empty).
func (g *PuzzleGame) IsSolved() bool {
   for rowIdx := 0; rowIdx < MaxPyramidRows; rowIdx++ {
       for colIdx := 0; colIdx < utils.PyramidRowSizes[rowIdx]; colIdx++ {
           if g.pyramid[rowIdx][colIdx] != -1 {
               return false
           }
       }
   }
   return true
}


// CalculateScore calculates the current score.
func (g *PuzzleGame) CalculateScore() int {
   matchingScore := g.matches * 50


   stonesRemainingScore := 0
   if g.IsSolved() {
       stonesRemaining := 0
       if g.hold != -1 {
           stonesRemaining++
       }
       for _, segment := range g.drawPile {
           stonesRemaining += len(segment)
       }
       stonesRemainingScore = stonesRemaining * 50
   }


   redrawCost := g.redraws * -50
   completionBonus := 0
   if g.IsSolved() {
       completionBonus = 500
   }


   completionPercentage := g.calculateCompletionPercentage()
   timeBonus := int(math.Floor(float64(g.timeRemaining) * completionPercentage * 6))


   totalScore := matchingScore + stonesRemainingScore + redrawCost + g.streakBonus + completionBonus + timeBonus
   return int(math.Max(0, float64(totalScore)))
}


// calculateCompletionPercentage calculates the percentage of the pyramid cleared.
func (g *PuzzleGame) calculateCompletionPercentage() float64 {
   clearedPositions := 0
   for rowIdx := 0; rowIdx < MaxPyramidRows; rowIdx++ {
       for colIdx := 0; colIdx < utils.PyramidRowSizes[rowIdx]; colIdx++ {
           if g.pyramid[rowIdx][colIdx] == -1 {
               clearedPositions++
           }
       }
   }
   return float64(clearedPositions) / TotalPyramidStones
}


// PrintState prints the current game state.
func (g *PuzzleGame) PrintState() {
   fmt.Println("\nPyramid:")
   for rowIdx := MaxPyramidRows - 1; rowIdx >= 0; rowIdx-- { // Iterate from G (6) down to A (0)
       fmt.Print(strings.Repeat("  ", rowIdx)) // Indentation
       for colIdx := 0; colIdx < utils.PyramidRowSizes[rowIdx]; colIdx++ {
           stone := g.pyramid[rowIdx][colIdx]
           if stone != -1 {
               fmt.Printf("%2d ", stone)
           } else {
               fmt.Print("-- ")
           }
       }
       fmt.Println()
   }


   fmt.Println("\nHold:", g.GetHoldValue())


   fmt.Println("\nDraw Pile:")
   for i := 0; i < g.numActiveSegments; i++ {
       fmt.Printf("Segment %d: ", i+1)
       for _, stone := range g.drawPile[i] {
           fmt.Printf("%2d ", stone)
       }
       if i == g.currentSegment {
           fmt.Print("(current)")
           if len(g.drawPile[i]) > 0 {
               fmt.Printf(" - DRW1: %d", g.drawPile[i][len(g.drawPile[i])-1])
           }
       }
       fmt.Println()
   }


   fmt.Printf("\nMatches: %d\n", g.matches)
   fmt.Printf("Streak: %d\n", g.streak)
   fmt.Printf("Streak Bonus: %d\n", g.streakBonus)
   fmt.Printf("Redraws: %d\n", g.redraws)
   fmt.Printf("Current Score: %d\n", g.CalculateScore())
   fmt.Printf("Solved: %t\n", g.IsSolved())
}


// GetHoldValue returns the value of the hold area, or "-" if empty.
func (g *PuzzleGame) GetHoldValue() interface{} {
   if g.hold != -1 {
       return g.hold
   }
   return "-"
}


// DeepCopy creates a new PuzzleGame instance with the same state as the original.
// This is crucial for running independent simulations in parallel.
func (g *PuzzleGame) DeepCopy() *PuzzleGame {
	newGame := NewPuzzleGame() // Start with a fresh game object

	// Copy simple fields
	newGame.hold = g.hold
	newGame.currentSegment = g.currentSegment
	newGame.matches = g.matches
	newGame.streak = g.streak
	newGame.streakBonus = g.streakBonus
	newGame.redraws = g.redraws
	newGame.timeRemaining = g.timeRemaining
	newGame.numActiveSegments = g.numActiveSegments

	// Deep copy the pyramid (array of arrays)
	for i := range g.pyramid {
		copy(newGame.pyramid[i][:], g.pyramid[i][:])
	}

	// Deep copy the draw pile (array of slices)
	for i := range g.drawPile {
		if g.drawPile[i] != nil {
			newGame.drawPile[i] = make([]int, len(g.drawPile[i]))
			copy(newGame.drawPile[i], g.drawPile[i])
		}
	}

	// Deep copy the moves slice
	if g.moves != nil {
		newGame.moves = make([]Move, len(g.moves))
		copy(newGame.moves, g.moves)
	}

	return newGame
}

// Reset restores the game state to match an original template game.
// This is much more efficient than DeepCopy for repeated simulations as it reuses memory.
func (g *PuzzleGame) Reset(original *PuzzleGame) {
	// Reset scalar fields
	g.hold = original.hold
	g.currentSegment = original.currentSegment
	g.matches = original.matches
	g.streak = original.streak
	g.streakBonus = original.streakBonus
	g.redraws = original.redraws
	g.timeRemaining = original.timeRemaining
	g.numActiveSegments = original.numActiveSegments

	// Reset the moves slice
	g.moves = g.moves[:0] // Efficiently clear the slice while retaining capacity

	// Copy the pyramid state
	for i := range original.pyramid {
		copy(g.pyramid[i][:], original.pyramid[i][:])
	}

	// Copy the draw pile state
	for i := range original.drawPile {
		// Ensure the destination slice has enough capacity
		if cap(g.drawPile[i]) < len(original.drawPile[i]) {
			g.drawPile[i] = make([]int, len(original.drawPile[i]))
		} else {
			g.drawPile[i] = g.drawPile[i][:len(original.drawPile[i])]
		}
		copy(g.drawPile[i], original.drawPile[i])
	}
}


// --- Public accessors for solver ---
func (g *PuzzleGame) PyramidValue(row, col int) int {
   return g.pyramid[row][col]
}
func (g *PuzzleGame) HoldValue() int {
   return g.hold
}
func (g *PuzzleGame) Redraws() int {
   return g.redraws
}
func (g *PuzzleGame) NumActiveSegments() int {
   return g.numActiveSegments
}


// DrawPile returns the entire drawPile array (array of slices).
func (g *PuzzleGame) DrawPile() [MaxDrawPileSegments][]int {
   return g.drawPile
}


