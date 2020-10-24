package generator

import (
	"math/rand"
	solver "sudoku-partner/solver"
	"time"
)

func GenRandomSudoku() string {
	const numDigit = 4
	dl := GenRandomSudokuWithNDigit(numDigit, 0)
	for !dl.Solve() {
		dl = GenRandomSudokuWithNDigit(numDigit, 0)
	}

	return dl.String()
}

func GenRandomSudokuWithNDigit(nDigit int, seed int64) solver.DancingLink {
	var row int
	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	r := rand.New(rand.NewSource(seed))
	dl := solver.NewDancingLink()
	for i := 0; i < nDigit; i++ {
		for row = r.Int() % 729; !dl.ContainsRow(row); {
			row = r.Int() % 729
		}
		dl.Insert(row/9, row%9+1)
	}
	return dl
}
