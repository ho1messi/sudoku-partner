package solver

import (
	"fmt"
	"strings"

	mapset "github.com/deckarep/golang-set"
)

type DancingLinkArray struct {
	condMatrix [729][324]bool
	rowFlags   [729]bool
	colFlags   [324]bool
	steps      mapset.Set
	board      [81]int
	deleteRows map[int][]int
	deleteCols map[int][]int
	colCount   int
}

func NewDancingLinkArray() DancingLinkArray {
	dl := DancingLinkArray{
		[729][324]bool{},
		[729]bool{},
		[324]bool{},
		mapset.NewSet(),
		[81]int{},
		make(map[int][]int),
		make(map[int][]int),
		324,
	}
	dl.initCondMatrix()
	return dl
}

func NewDancingLinkArrayFromString(sudoku string) DancingLinkArray {
	dl := NewDancingLinkArray()
	for i := 0; i < 81; i++ {
		if sudoku[i] >= '1' && sudoku[i] <= '9' {
			dl.Insert(i/9, i%9, int(sudoku[i]-'0'))
		}
	}
	return dl
}

func (dl *DancingLinkArray) Insert(row int, col int, digit int) {
	dl.board[row*9+col] = digit
	subgrid := row/3*3 + col/3
	cols := []int{
		row*9 + col,
		row*9 + digit + 80,
		col*9 + digit + 161,
		subgrid*9 + digit + 242,
	}
	for _, col := range cols {
		dl.removeCol(col, -1)
	}
}

func (dl *DancingLinkArray) Solve() {
	dl.solveR()
}

func (dl *DancingLinkArray) Info() string {
	rowCount, colCount := 0, 0
	for _, rowFlag := range dl.rowFlags {
		if !rowFlag {
			rowCount++
		}
	}
	for _, colFlag := range dl.colFlags {
		if !colFlag {
			colCount++
		}
	}
	return fmt.Sprintf("row:%d col:%d %d", rowCount, colCount, dl.colCount)
}

func (dl DancingLinkArray) String() string {
	strs := []string{}
	for _, digit := range dl.board {
		strs = append(strs, digitStrs[digit])
	}
	return strings.Join(strs, "")
}

func (dl *DancingLinkArray) initCondMatrix() {
	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
			for digit := 1; digit < 10; digit++ {
				r, subgrid := row*81+col*9+digit-1, row/3*3+col/3
				dl.condMatrix[r][row*9+col] = true
				dl.condMatrix[r][row*9+digit+80] = true
				dl.condMatrix[r][col*9+digit+161] = true
				dl.condMatrix[r][subgrid*9+digit+242] = true
			}
		}
	}
}

func (dl *DancingLinkArray) solveR() bool {
	cols := dl.sortCols()
	for _, col := range cols {
		for row := 0; row < len(dl.rowFlags); row++ {
			if !dl.rowFlags[row] && dl.condMatrix[row][col] {
				dl.selectRow(row)
				if dl.colCount == 0 || dl.solveR() {
					return true
				}
				dl.unselectRow(row)
			}
		}
	}
	return false
}

func (dl *DancingLinkArray) sortCols() []int {
	colOrder, colSize := make([]int, 324), make([]int, 324)
	nCols := 0
	for col := range dl.colFlags {
		if !dl.colFlags[col] {
			i, count := 0, 0
			for r := 0; r < len(dl.rowFlags); r++ {
				if !dl.rowFlags[r] && dl.condMatrix[r][col] {
					count++
				}
			}
			for ; i < nCols && colSize[i] <= count; i++ {
			}
			for j := nCols; j >= i; j-- {
				colOrder[j+1] = colOrder[j]
				colSize[j+1] = colSize[j]
			}
			colOrder[i] = col
			colSize[i] = count
			nCols++
		}
	}
	return colOrder[:nCols]
}

func (dl *DancingLinkArray) selectRow(row int) {
	dl.steps.Add(row)
	index, digit := rowToBoard(row)
	dl.board[index] = digit
	for col, colFlag := range dl.colFlags {
		if dl.condMatrix[row][col] && !colFlag {
			dl.removeCol(col, row)
		}
	}
}

func (dl *DancingLinkArray) unselectRow(row int) {
	dl.steps.Remove(row)
	index, _ := rowToBoard(row)
	dl.board[index] = 0
	dl.rowFlags[row] = false
	for _, col := range dl.deleteCols[row] {
		dl.colFlags[col] = false
		dl.colCount++
	}
	for _, r := range dl.deleteRows[row] {
		dl.rowFlags[r] = false
	}
	delete(dl.deleteCols, row)
	delete(dl.deleteRows, row)
}

func (dl *DancingLinkArray) removeCol(col int, byRow int) {
	dl.deleteCols[byRow] = append(dl.deleteCols[byRow], col)
	dl.colFlags[col] = true
	dl.colCount--
	for row, rowFlag := range dl.rowFlags {
		if row != byRow && !rowFlag && dl.condMatrix[row][col] {
			dl.deleteRows[byRow] = append(dl.deleteRows[byRow], row)
			dl.rowFlags[row] = true
		}
	}
}

func rowToBoard(row int) (int, int) {
	return row / 9, row%9 + 1
}
