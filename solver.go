package main

import "time"

// ValidRow searches row of board for val
// result is false if val is found
func ValidRow(row, val int, board *[][]int) bool {
	for i := 0; i < 9; i++ {
		if val == (*board)[row][i] {
			return false
		}
	}
	return true
}

// ValidCol searches col of board for val
// result is false if val is found
func ValidCol(col, val int, board *[][]int) bool {
	for i := 0; i < 9; i++ {
		if val == (*board)[i][col] {
			return false
		}
	}
	return true
}

// ValidSquare searches 3x3 square of board that contains row/col for val
// result is false if val is found
func ValidSquare(row, col, val int, board *[][]int) bool {
	// set row and col to starting index of square
	row = row / 3 * 3
	col = col / 3 * 3
	for i := row; i < row+3; i++ {
		for j := col; j < col+3; j++ {
			if val == (*board)[i][j] {
				return false
			}
		}
	}
	return true
}

// ValidPlacement uses goroutines to check for valid square, row, and column simutaneously
// returns false if a value already exists
func ValidPlacement(row, col, val int, board *[][]int) bool {
	if (*board)[row][col] != 0 { // space must be empty
		return false
	}
	return ValidSquare(row, col, val, board) && ValidRow(row, val, board) && ValidCol(col, val, board)

}

// FindEmpty searches the board for the first empty space
// returns the index of first found
func FindEmpty(board *[][]int) (bool, int, int) {
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if (*board)[i][j] == 0 { // empty space found
				return true, i, j
			}
		}
	}
	return false, -1, -1
}

// SolveBoard recursively solve board and backtrack when necessary
// tests digits 1-9 for first empty space found
// when current solution becomes invalid, backtrack to last valid solution
func SolveBoard(board *[][]int) bool {
	for !solving {
		continue
	}
	// search for empty space and set row, col as the index
	empty, row, col := FindEmpty(board)
	if !empty { // puzzle is finished
		solved = true
		endTime = time.Since(startTime)
		return true
	}
	// consider digits 1-9
	for n := 1; n <= 9; n++ {
		if ValidPlacement(row, col, n, board) {
			(*board)[row][col] = n
			if SolveBoard(board) {
				return true
			}
			// path did not lead to solution, reverse progress
			(*board)[row][col] = 0
		}
	}
	// trigger backtrack
	return false
}

// CreateBoard initializes a blank 9x9 sudoku board
func CreateBoard() [][]int {
	board := make([][]int, 9)
	for i := 0; i < 9; i++ {
		board[i] = make([]int, 9)
	}
	return board
}
