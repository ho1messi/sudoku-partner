package solver

import (
	"testing"

	mapset "github.com/deckarep/golang-set"
)

func TestDancingLinkInit(t *testing.T) {
	dl := NewDancingLink()
	if dl.colNodes[0].left != dl.head || dl.colNodes[323].right != dl.head {
		t.Errorf("\ncolNode error.\n col: 0 or 323")
	}
	for i, colNode := range dl.colNodes {
		if (i > 0 && colNode.left != dl.colNodes[i-1]) ||
			(i < 323 && colNode.right != dl.colNodes[i+1]) {
			t.Errorf("\ncolNode error.\ncol: %d", i)
		}
		count := countColSize(&dl, i)
		if count != 9 {
			t.Errorf("\nsize of col is not 9!\ncol: %d, size: %d", i, count)
		}
	}

	for i := range dl.rowNodes {
		count := countRowSize(&dl, i)
		if count != 4 {
			t.Errorf("\nsize of row is not 4!\nrow: %d, size: %d", i, count)
		}
	}
}

func TestDancingLinkInsert(t *testing.T) {
	for _, puzzle := range easyPuzzles {
		sudoku := puzzle[0]
		dl := NewDancingLink()
		for i := 0; i < 81; i++ {
			if sudoku[i] >= '1' && sudoku[i] <= '9' {
				cols, rows, rowCols := getRelativeRowAndCols(&dl, i, int(sudoku[i]-'0'))
				dl.Insert(i, int(sudoku[i]-'0'))

				for colNode := dl.head.right; colNode != dl.head; colNode = colNode.right {
					if cols.Contains(colNode.col) {
						t.Errorf("\nInsert error!\nfailed to remove col: %d", colNode.col)
					}
				}
				for _, col := range rowCols.ToSlice() {
					colNode := dl.colNodes[col.(int)]
					for node := colNode.down; node != colNode; node = node.down {
						if rows.Contains(node.row) {
							t.Errorf("\nInsert error!\nfailed to remove row: %d", node.row)
						}
					}
				}

			}
		}

		result := dl.String()
		if result != sudoku {
			t.Errorf("\nCannot insert correct digit\nsudoku: %s\nresult: %s", sudoku, result)
		}
	}
}

func TestDancingLinkSolve(t *testing.T) {
	for _, puzzle := range veryHardPuzzles {
		sudoku, result := puzzle[0], puzzle[1]
		dl := NewDancingLinkFromString(sudoku)
		dl.Solve()
		solved := dl.String()
		if solved != result {
			t.Errorf("\nError case:\nsudoku: %s\nsolved: %s\nresult: %s", sudoku, solved, result)
		}
	}
}

func countColSize(dl *DancingLink, col int) int {
	count := 0
	colNode := dl.colNodes[col]
	for node := colNode.down; node != colNode; node = node.down {
		count++
	}
	return count
}

func countRowSize(dl *DancingLink, row int) int {
	count := 0
	rowNode := dl.rowNodes[row]
	for node := rowNode.right; node != rowNode; node = node.right {
		count++
	}
	return count
}

func getRelativeRowAndCols(dl *DancingLink, index int, digit int) (mapset.Set, mapset.Set, mapset.Set) {
	row := index*9 + digit - 1
	rowNode := dl.rowNodes[row]
	cols, rows, rowCols := mapset.NewSet(), mapset.NewSet(), mapset.NewSet()
	for node := rowNode.right; node != rowNode; node = node.right {
		cols.Add(node.col)
	}
	for _, col := range cols.ToSlice() {
		colNode := dl.colNodes[col.(int)]
		for node1 := colNode.down; node1 != colNode; node1 = node1.down {
			rows.Add(node1.row)
			rNode := dl.rowNodes[node1.row]
			for node2 := rNode.right; node2 != rNode; node2 = node2.right {
				rowCols.Add(node2.col)
			}
		}
	}
	rowCols = rowCols.Difference(cols)
	return cols, rows, rowCols
}
