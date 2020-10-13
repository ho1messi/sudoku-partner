package solver

import (
	"fmt"
	"strings"

	mapset "github.com/deckarep/golang-set"
)

var digitStrs = []string{
	".", "1", "2", "3", "4", "5", "6", "7", "8", "9",
}

type DancingLink1 struct {
	condMatrix [729][324]bool
	rowFlags   [729]bool
	colFlags   [324]bool
	steps      mapset.Set
	board      [81]int
	deleteRows map[int][]int
	deleteCols map[int][]int
	colCount   int
}

func NewDancingLink1() DancingLink1 {
	dl := DancingLink1{
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

func NewDancingLink1FromString(sudoku string) DancingLink1 {
	dl := NewDancingLink1()
	for i := 0; i < 81; i++ {
		if sudoku[i] >= '1' && sudoku[i] <= '9' {
			dl.Insert(i/9, i%9, int(sudoku[i]-'0'))
		}
	}
	return dl
}

func (dl *DancingLink1) Insert(row int, col int, digit int) {
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

func (dl *DancingLink1) Solve() {
	dl.solveR(0)
}

func (dl *DancingLink1) Info() string {
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

func (dl DancingLink1) String() string {
	strs := []string{}
	for _, digit := range dl.board {
		strs = append(strs, digitStrs[digit])
	}
	return strings.Join(strs, "")
	// str := ""
	// for _, digit := range dl.board {
	// 	str += digitStrs[digit]
	// }
	// return str
}

func (dl *DancingLink1) initCondMatrix() {
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

func (dl *DancingLink1) solveR(startCol int) {
	for col := startCol; col < len(dl.colFlags); col++ {
		if !dl.colFlags[col] {
			for row := 0; row < len(dl.rowFlags); row++ {
				if !dl.rowFlags[row] && dl.condMatrix[row][col] {
					dl.selectRow(row)
					if dl.colCount == 0 {
						return
					}
					dl.solveR(col + 1)
					dl.unselectRow(row)
				}
			}
		}
	}
}

func (dl *DancingLink1) selectRow(row int) {
	dl.steps.Add(row)
	index, digit := rowToBoard(row)
	dl.board[index] = digit
	for col, colFlag := range dl.colFlags {
		if dl.condMatrix[row][col] && !colFlag {
			dl.removeCol(col, row)
			// dl.deleteCols[row] = append(dl.deleteCols[row], col)
			// dl.colFlags[col] = true
			// dl.colCount--
			// for r := 0; r < len(dl.rowFlags); r++ {
			// 	if r != row && dl.rowFlags[r] && dl.condMatrix[r][col] {
			// 		dl.deleteRows[row] = append(dl.deleteRows[row], r)
			// 		dl.rowFlags[r] = true
			// 	}
			// }
		}
	}
}

func (dl *DancingLink1) unselectRow(row int) {
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

func (dl *DancingLink1) removeCol(col int, byRow int) {
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
