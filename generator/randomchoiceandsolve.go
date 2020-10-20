package generator

import (
	"math/rand"
	solver "sudoku-partner/solver"
	"time"
)

func GenRandomSudoku() string {
	const numDigit = 12
	dl := genRandomSudokuWithNDigit(numDigit)
	for !dl.Solve() {
		dl = genRandomSudokuWithNDigit(numDigit)
	}

	return dl.String()
}

func genRandomSudokuWithNDigit(nDigit int) solver.DancingLink {
	var row int
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	dl := solver.NewDancingLink()
	for i := 0; i < nDigit; i++ {
		for row = r.Int() % 729; !dl.ContainsRow(row); {
			row = r.Int() % 729
		}
		dl.Insert(row/9, row%9+1)
	}
	return dl
}
