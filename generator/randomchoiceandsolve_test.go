package generator

import (
	"testing"

	mapset "github.com/deckarep/golang-set"
)

func TestGenRandomSudoku(t *testing.T) {
	for i := 0; i < 1000; i++ {
		sudoku := GenRandomSudoku()
		if !isSudokuValid(sudoku) {
			t.Errorf("\n%s", sudoku)
		}
	}
}

func BenchmarkGenRandomSudokuWithNDigit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenRandomSudokuWithNDigit(11, 0)
	}
}

func BenchmarkGenRandomSudoku(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenRandomSudoku()
	}
}

func isSudokuValid(sudoku string) bool {
	digits := make([]int, 81)
	for i := 0; i < 81; i++ {
		if sudoku[i] < '1' || sudoku[i] > '9' {
			return false
		}
		digits[i] = int(sudoku[i] - '0')
	}
	for i := 0; i < 9; i++ {
		sets := [3]mapset.Set{mapset.NewSet(), mapset.NewSet(), mapset.NewSet()}
		for j := 0; j < 9; j++ {
			indices := []int{
				i*9 + j,
				j*9 + i,
				(i/3*3+j/3)*9 + i%3*3 + j%3,
			}
			for k := 0; k < 3; k++ {
				sets[k].Add(digits[indices[k]])
			}
		}
		for k := 0; k < 3; k++ {
			if sets[k].Cardinality() != 9 {
				return false
			}
		}
	}
	return true
}
