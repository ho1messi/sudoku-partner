package solver

import "testing"

func TestBruteForceInsert(t *testing.T) {
	for _, puzzle := range easyPuzzles {
		sudoku := puzzle[0]
		bf := NewBoardBruteForceFromString(sudoku)
		result := bf.String()
		if result != sudoku {
			t.Errorf("\nError case:\nsudoku: %s\nresult: %s\n", sudoku, result)
		}
	}
}

func TestBruteForceSolve(t *testing.T) {
	for _, puzzle := range easyPuzzles {
		sudoku, result := puzzle[0], puzzle[1]
		bf := NewBoardBruteForceFromString(sudoku)
		info := bf.Info()
		bf.Solve()
		solved := bf.String()
		if solved != result {
			t.Errorf("\nError case:\nsudoku: %s\nsolved: %s\nresult: %s\ninfo: %s", sudoku, solved, result, info)
		}
	}
}
