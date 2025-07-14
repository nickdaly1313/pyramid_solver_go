package utils

import (
	"fmt"
	"math/rand" // Used for ShuffleArray
	"strconv"
	"strings"
)

// PyramidRowSizes maps row indices (0 for A, 6 for G) to their number of columns.
// This is used for validation and internal representation.
var PyramidRowSizes = []int{7, 6, 5, 4, 3, 2, 1} // Index 0=A, 1=B, ..., 6=G

// RowCharToIndex converts a row character ('A'-'G') to its 0-indexed integer.
func RowCharToIndex(r rune) int {
	return int(r - 'A')
}

// ColCharToIndex converts a column character ('a'-'g') to its 0-indexed integer.
func ColCharToIndex(c rune) int {
	return int(c - 'a')
}

// IndexToRowChar converts a 0-indexed row integer to its character ('A'-'G').
func IndexToRowChar(idx int) rune {
	return rune('A' + idx)
}

// IndexToColChar converts a 0-indexed column integer to its character ('a'-'g').
func IndexToColChar(idx int) rune {
	return rune('a' + idx)
}

// StringToIndices converts a position string (e.g., "A1") to row and column indices.
// Returns row_idx, col_idx, and an error if the format is invalid.
func StringToIndices(posStr string) (int, int, error) {
	if len(posStr) < 2 || len(posStr) > 3 { // Handle A1-A7, B1-B6, etc. (e.g., A1, A10 - though not used here)
		return -1, -1, fmt.Errorf("invalid position string length: %s", posStr)
	}

	rowChar := rune(posStr[0])
	colStr := posStr[1:] // Extract column part (can be 1 or 2 digits)

	if rowChar < 'A' || rowChar > 'G' {
		return -1, -1, fmt.Errorf("invalid row character in position string: %s", posStr)
	}

	rowIdx := RowCharToIndex(rowChar)
	colIdx, err := strconv.Atoi(colStr)
	if err != nil {
		return -1, -1, fmt.Errorf("invalid column number in position string: %s", posStr)
	}
	colIdx-- // Adjust to 0-based index

	if rowIdx < 0 || rowIdx >= len(PyramidRowSizes) || colIdx < 0 || colIdx >= PyramidRowSizes[rowIdx] {
		return -1, -1, fmt.Errorf("position %s out of pyramid bounds", posStr)
	}

	return rowIdx, colIdx, nil
}

// IndicesToString converts row and column indices to a position string (e.g., "A1").
// Returns the string and an error if indices are out of bounds.
func IndicesToString(rowIdx, colIdx int) (string, error) {
	if rowIdx < 0 || rowIdx >= len(PyramidRowSizes) {
		return "", fmt.Errorf("row index out of bounds: %d", rowIdx)
	}
	rowChar := IndexToRowChar(rowIdx)

	if colIdx < 0 || colIdx >= PyramidRowSizes[rowIdx] {
		return "", fmt.Errorf("column index %d out of bounds for row %c", colIdx, rowChar)
	}

	return fmt.Sprintf("%c%d", rowChar, colIdx+1), nil // Add 1 for 1-based column
}

// ParseInts parses a space-separated string of integers.
func ParseInts(s string, expectedCount int) ([]int, error) {
	parts := strings.Fields(s)
	if len(parts) != expectedCount {
		return nil, fmt.Errorf("expected %d values, got %d", expectedCount, len(parts))
	}
	nums := make([]int, expectedCount)
	for i, p := range parts {
		n, err := strconv.Atoi(p)
		if err != nil {
			return nil, fmt.Errorf("invalid number '%s': %w", p, err)
		}
		if n < 1 || n > 13 {
			return nil, fmt.Errorf("value %d out of range (1-13)", n)
		}
		nums[i] = n
	}
	return nums, nil
}

// ShuffleArray shuffles a slice of integers.
func ShuffleArray(arr []int) []int {
	shuffled := make([]int, len(arr))
	copy(shuffled, arr)
	for i := len(shuffled) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}
	return shuffled
}
