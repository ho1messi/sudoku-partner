package main

import (
	"fmt"
	"sudoku-partner/solver"
)

func main() {
	// sudoku := os.Args[1]
	// sudoku := ".6.593...9.1...5...3.4...9.1.8.2...44..3.9..12...1.6.9.8...6.2...4...8.7...785.1."
	sudoku := "981..3.4.6...7925.2751.698319.4.75325.831.7.97236.541831.7.4.9..6923.1.5.579.1324"
	dl := solver.NewDancingLink1FromString(sudoku)
	fmt.Println(dl)
	fmt.Println(dl.Info())
	dl.Solve()
	fmt.Println(dl)
}
