package solver

import (
	"strings"

	mapset "github.com/deckarep/golang-set"
)

type DancingLink struct {
	board    [81]int
	head     *DancingLinkNode
	colNodes [](*DancingLinkNode)
	rowNodes [](*DancingLinkNode)
	colSet   mapset.Set
	rowSet   mapset.Set
}

type DancingLinkNode struct {
	left  *DancingLinkNode
	right *DancingLinkNode
	up    *DancingLinkNode
	down  *DancingLinkNode
	size  int
	row   int
	col   int
}

func NewDancingLink() DancingLink {
	var dl DancingLink
	dl.initDancingLinkRows()
	return dl
}

func NewDancingLinkFromString(sudoku string) DancingLink {
	dl := NewDancingLink()
	for i := 0; i < 81; i++ {
		if sudoku[i] >= '1' && sudoku[i] <= '9' {
			dl.Insert(i, int(sudoku[i]-'0'))
		}
	}
	return dl
}

func (dl *DancingLink) Insert(index int, digit int) {
	dl.board[index] = digit
	row := index*9 + digit - 1
	rowNode := dl.rowNodes[row]
	for node := rowNode.right; node != rowNode; node = node.right {
		dl.coverCol(node.col)
	}
}

func (dl *DancingLink) Solve() bool {
	dl.solve(nil)
	return dl.head.right == dl.head
}

func (dl *DancingLink) GetAllResult() []string {
	if dl.head.right == dl.head {
		return []string{dl.String()}
	} else {
		var results []string
		dl.solve(&results)
		return results
	}
}

func (dl *DancingLink) ContainsCol(col int) bool {
	return dl.colSet.Contains(col)
}

func (dl *DancingLink) ContainsRow(row int) bool {
	return dl.rowSet.Contains(row)
}

func (dl *DancingLink) IsSolved() bool {
	return dl.head.right == dl.head
}

func (dl DancingLink) String() string {
	strs := []string{}
	for _, digit := range dl.board {
		strs = append(strs, digitStrs[digit])
	}
	return strings.Join(strs, "")
}

func (dl *DancingLink) initDancingLinkRows() {
	dl.head = newDancingLinkNode()
	dl.colNodes, dl.colSet = genColNodes(dl.head, 324)
	dl.rowNodes, dl.rowSet = make([](*DancingLinkNode), 729), mapset.NewSet()
	for i := 0; i < 81; i++ {
		for digit := 1; digit < 10; digit++ {
			index, subgrid := i*9+digit-1, i/9/3*3+i%9/3
			cols := []int{
				i/9*9 + i%9,
				i/9*9 + digit + 80,
				i%9*9 + digit + 161,
				subgrid*9 + digit + 242,
			}
			dl.rowNodes[index] = genRowNodes(dl.colNodes, index, cols)
			dl.rowSet.Add(i*9 + digit - 1)
		}
	}
}

func genColNodes(head *DancingLinkNode, nCols int) ([](*DancingLinkNode), mapset.Set) {
	colNodes := make([](*DancingLinkNode), nCols)
	colSet := mapset.NewSet()
	for i := 0; i < nCols; i++ {
		node := newDancingLinkNode()
		node.right = head
		node.left = head.left
		head.left.right = node
		head.left = node
		node.col = i
		node.row = -1
		colNodes[i] = node
		colSet.Add(i)
	}
	return colNodes, colSet
}

func genRowNodes(colNodes [](*DancingLinkNode), row int, cols []int) *DancingLinkNode {
	rowNode := newDancingLinkNode()
	rowNode.row = row
	rowNode.col = -1
	for _, i := range cols {
		col, node := colNodes[i], newDancingLinkNode()
		node.col = i
		node.row = row
		node.up = col.up
		node.down = col
		col.up.down = node
		col.up = node
		col.size++

		node.left = rowNode.left
		node.right = rowNode
		rowNode.left.right = node
		rowNode.left = node
	}
	return rowNode
}

func (dl *DancingLink) coverCol(col int) {
	colNode := dl.colNodes[col]
	colNode.left.right = colNode.right
	colNode.right.left = colNode.left
	dl.colSet.Remove(col)
	for node1 := colNode.down; node1 != colNode; node1 = node1.down {
		rowNode := dl.rowNodes[node1.row]
		dl.rowSet.Remove(node1.row)
		for node2 := rowNode.right; node2 != rowNode; node2 = node2.right {
			if node2 != node1 {
				node2.up.down = node2.down
				node2.down.up = node2.up
				dl.colNodes[node2.col].size--
			}
		}
	}
}

func (dl *DancingLink) uncoverCol(col int) {
	colNode := dl.colNodes[col]
	colNode.left.right = colNode
	colNode.right.left = colNode
	dl.colSet.Add(col)
	for node1 := colNode.up; node1 != colNode; node1 = node1.up {
		rowNode := dl.rowNodes[node1.row]
		dl.rowSet.Add(node1.row)
		for node2 := rowNode.left; node2 != rowNode; node2 = node2.left {
			if node2 != node1 {
				node2.up.down = node2
				node2.down.up = node2
				dl.colNodes[node2.col].size++
			}
		}
	}
}

func (dl *DancingLink) coverRow(row int) {
	rowNode := dl.rowNodes[row]
	for node := rowNode.right; node != rowNode; node = node.right {
		dl.coverCol(node.col)
	}
}

func (dl *DancingLink) uncoverRow(row int) {
	rowNode := dl.rowNodes[row]
	for node := rowNode.left; node != rowNode; node = node.left {
		dl.uncoverCol(node.col)
	}
}

func (dl *DancingLink) addStep(row int) {
	dl.board[row/9] = row%9 + 1
}

func (dl *DancingLink) removeStep(row int) {
	dl.board[row/9] = 0
}

func (dl *DancingLink) solve(results *[]string) {
	if dl.head.right != dl.head {
		colNode := dl.selectCol()
		for node := colNode.down; node != colNode; node = node.down {
			dl.addStep(node.row)
			dl.coverRow(node.row)

			dl.solve(results)
			if dl.head.right == dl.head {
				if results != nil {
					*results = append(*results, dl.String())
				} else {
					return
				}
			}

			dl.uncoverRow(node.row)
			dl.removeStep(node.row)
		}
	}
}

func (dl *DancingLink) selectCol() *DancingLinkNode {
	var colNode *DancingLinkNode
	for node := dl.head.right; node != dl.head; node = node.right {
		if colNode == nil || node.size < colNode.size {
			colNode = node
		}
	}
	return colNode
}

func newDancingLinkNode() *DancingLinkNode {
	node := new(DancingLinkNode)
	node.up = node
	node.down = node
	node.left = node
	node.right = node
	return node
}
