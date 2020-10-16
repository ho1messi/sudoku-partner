package solver

import (
	"fmt"
	"strings"
)

type BoardBruteForce struct {
	board     [81]int
	units     [27][9]int
	cellUnits [81][3]int
	stepToGo  int
}

func NewBoardBruteForce() BoardBruteForce {
	bf := BoardBruteForce{
		[81]int{},
		[27][9]int{},
		[81][3]int{},
		81,
	}
	bf.initCellUnits()
	return bf
}

func NewBoardBruteForceFromString(sudoku string) BoardBruteForce {
	bf := NewBoardBruteForce()
	for i := 0; i < 81; i++ {
		if sudoku[i] >= '1' && sudoku[i] <= '9' {
			bf.Insert(i/9, i%9, int(sudoku[i]-'0'))
		}
	}
	return bf
}

func (bf *BoardBruteForce) Insert(row int, col int, digit int) {
	index := row*9 + col
	if bf.board[index] == 0 {
		bf.stepToGo--
	}
	bf.board[index] = digit
}

func (bf *BoardBruteForce) Remove(row int, col int) {
	index := row*9 + col
	if bf.board[index] > 0 {
		bf.stepToGo++
	}
	bf.board[index] = 0
}

func (bf *BoardBruteForce) Solve() {
	bf.solveR(0)
}

func (bf BoardBruteForce) String() string {
	strs := []string{}
	for _, digit := range bf.board {
		strs = append(strs, digitStrs[digit])
	}
	return strings.Join(strs, "")
}

func (bf *BoardBruteForce) Info() string {
	return fmt.Sprintf("Step to go: %d", bf.stepToGo)
}

func (bf *BoardBruteForce) initCellUnits() {
	for i := range bf.board {
		row, col := i/9, i%9
		subgrid := row/3*3 + col/3
		cellUnits := []int{row, col + 9, subgrid + 18}
		unitOffsets := []int{col, row, row%3*3 + col%3}
		for j := range cellUnits {
			bf.cellUnits[i][j] = cellUnits[j]
			bf.units[cellUnits[j]][unitOffsets[j]] = i
		}
	}
}

func (bf *BoardBruteForce) solveR(startIndex int) bool {
	for i, digit := range bf.board {
		if i < startIndex || digit > 0 {
			continue
		}
		for d := 1; d < 10; d++ {
			conflictFlag := false
			for _, unit := range bf.cellUnits[i] {
				for _, ui := range bf.units[unit] {
					if bf.board[ui] == d {
						conflictFlag = true
						break
					}
				}
				if conflictFlag {
					break
				}
			}
			if !conflictFlag {
				bf.Insert(i/9, i%9, d)
				if bf.stepToGo == 0 || bf.solveR(i+1) {
					return true
				}
				bf.Remove(i/9, i%9)
				bf.board[i] = 0
			}
		}
	}
	return false
}
