package solver

import "testing"

func TestDancingLinkArraySolve(t *testing.T) {
	for _, puzzle := range easyPuzzles {
		sudoku, result := puzzle[0], puzzle[1]
		dl := NewDancingLinkArrayFromString(sudoku)
		info := dl.Info()
		dl.Solve()
		solved := dl.String()
		if solved != result {
			t.Errorf("\nError case:\nsudoku: %s\nsolved: %s\nresult: %s\ninfo: %s", sudoku, solved, result, info)
		}
	}
}
