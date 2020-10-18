package solver

import "strings"

type DancingLink struct {
	board    [81]int
	head     *DancingLinkNode
	colNodes [](*DancingLinkNode)
	rowNodes [](*DancingLinkNode)
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

func (dl *DancingLink) Solve() {
	if dl.head.right != dl.head {
		colNode := dl.head.right
		dl.coverCol(colNode.col)
		for node1 := colNode.down; node1 != colNode; node1 = node1.down {
			dl.addStep(node1.row)
			rowNode := dl.rowNodes[node1.row]
			for node2 := rowNode.right; node2 != rowNode; node2 = node2.right {
				if node2 != node1 {
					dl.coverCol(node2.col)
				}
			}
			dl.Solve()
			if dl.head.right == dl.head {
				return
			}
			for node2 := rowNode.left; node2 != rowNode; node2 = node2.left {
				if node2 != node1 {
					dl.uncoverCol(node2.col)
				}
			}
			dl.removeStep(node1.row)
		}
		dl.uncoverCol(colNode.col)
	}
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
	dl.colNodes = genColNodes(dl.head, 324)
	dl.rowNodes = make([](*DancingLinkNode), 729)
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
		}
	}
}

func genColNodes(head *DancingLinkNode, nCols int) [](*DancingLinkNode) {
	colNodes := make([](*DancingLinkNode), nCols)
	for i := 0; i < nCols; i++ {
		node := newDancingLinkNode()
		node.right = head
		node.left = head.left
		head.left.right = node
		head.left = node
		node.col = i
		node.row = -1
		colNodes[i] = node
	}
	return colNodes
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
	for node1 := colNode.down; node1 != colNode; node1 = node1.down {
		// dl.coverRow(node1.row)
		rowNode := dl.rowNodes[node1.row]
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
	for node1 := colNode.up; node1 != colNode; node1 = node1.up {
		rowNode := dl.rowNodes[node1.row]
		for node2 := rowNode.left; node2 != rowNode; node2 = node2.left {
			if node2 != node1 {
				node2.up.down = node2
				node2.down.up = node2
				dl.colNodes[node2.col].size++
			}
		}
	}
}

func (dl *DancingLink) addStep(row int) {
	dl.board[row/9] = row%9 + 1
}

func (dl *DancingLink) removeStep(row int) {
	dl.board[row/9] = 0
}

func newDancingLinkNode() *DancingLinkNode {
	node := new(DancingLinkNode)
	node.up = node
	node.down = node
	node.left = node
	node.right = node
	return node
}
