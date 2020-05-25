package main

import "github.com/amorydoerr/sudoku-solver/sudoku"

var sudokuBoard = [][]int{
	{8, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 3, 6, 0, 0, 0, 0, 0},
	{0, 7, 0, 0, 9, 0, 2, 0, 0},
	{0, 5, 0, 0, 0, 7, 0, 0, 0},
	{0, 0, 0, 0, 4, 5, 7, 0, 0},
	{0, 0, 0, 1, 0, 0, 0, 3, 0},
	{0, 0, 1, 0, 0, 0, 0, 6, 8},
	{0, 0, 8, 5, 0, 0, 0, 1, 0},
	{0, 9, 0, 0, 0, 0, 4, 0, 0},
}

// var sudokuBoard = CreateBoard()
var solving, solved, failed bool
var startButton = new(widget.Clickable)
var startTime time.Time
var endTime time.Duration


//
// BEGINNING OF GUI CODE
//



// driver for solving algorithm
func main() {
	go func () {
	failed = !SolveBoard(&sudokuBoard)
	}()
	CreateWindow()
}
