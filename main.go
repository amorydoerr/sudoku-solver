package main

import "fmt"

// ValidRow searches row of board for val
// outputs result to bool channel
func ValidRow(row, val int, board *[][]int, c chan<- bool) {
	for i := 0; i < 9; i++ {
		if val == (*board)[row][i] {
			c <- false
			return
		}
	}
	c <- true
}

// ValidCol searches col of board for val
// result is false if val is found
// outputs result to bool channel
func ValidCol(col, val int, board *[][]int, c chan<- bool) {
	for i := 0; i < 9; i++ {
		if val == (*board)[i][col] {
			c <- false
			return
		}
	}
	c <- true
}

// ValidSquare searches 3x3 square of board that contains row/col for val
// result is false if val is found
// outputs result to bool channel
func ValidSquare(row, col, val int, board *[][]int, c chan<- bool) {
	// set row and col to starting index of square
	row = row / 3 * 3
	col = col / 3 * 3
	for i := row; i < row+3; i++ {
		for j := col; j < col+3; j++ {
			if val == (*board)[i][j] {
				c <- false
				return
			}
		}
	}
	c <- true
}

// ValidPlacement uses goroutines to check for valid square, row, and column simutaneously
// returns false if a value already exists
func ValidPlacement(row, col, val int, board *[][]int) bool {
	if (*board)[row][col] != 0 { // space must be empty
		return false
	}
	// use go concurrency to perform all three checks simutaneously
	results := make(chan bool, 3)
	go ValidSquare(row, col, val, board, results)
	go ValidRow(row, val, board, results)
	go ValidCol(col, val, board, results)
	return <-results && <-results && <-results
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
	// search for empty space and set row, col as the index
	empty, row, col := FindEmpty(board)
	if !empty { // puzzle is finished
		PrintBoard(board)
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

// PrintBoard outputs the board with grid seperation
func PrintBoard(board *[][]int) {
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if j == 2 || j == 5 {
				fmt.Printf("%d|", (*board)[i][j])
			} else {
				fmt.Printf("%d ", (*board)[i][j])
			}
		}
		if i == 2 || i == 5 {
			fmt.Println("\n-----------------")
		} else {
			fmt.Println()
		}
	}
}

// driver for solving algorithm
func main() {
	testBoard := [][]int{
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
	fmt.Println()
	if !SolveBoard(&testBoard) {
		fmt.Println("No Solution")
	}
}
