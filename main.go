package main

import (
	"fmt"
	"log"
	"strconv"
	"time"
	// "image/color"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	// "gioui.org/f32"
	// "gioui.org/op/paint"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

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
var solving, solved bool
var startButton = new(widget.Clickable)
var startTime time.Time
var endTime time.Duration

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

// Start waits to start solving until button is pressed
// func Start(board *[][]int) {
// 	for {
// 		if solving {
// 			SolveBoard(board)
// 			break
// 		}
// 	}
// }

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

// UI CODE BEGINS HERE
// !!!!!!!!!!!!!!!!!!!!
// UI CODE BEGINS HERE

// CreateWindow opens a new gioui window
func CreateWindow() {
	gofont.Register()
	go func() {
		w := app.NewWindow(app.Size(unit.Dp(450), unit.Dp(500)), app.Title("Sudoku Solver"))
		if err := WindowLoop(w); err != nil {
			log.Fatal(err)
		}
	}()
	app.Main()
}

// WindowLoop controls the drawing of the widgets
func WindowLoop(w *app.Window) error {
	th := material.NewTheme()
	gtx := new(layout.Context)
	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx.Reset(e.Queue, e.Config, e.Size)

			DrawGrid(gtx, th)

			LayoutGrid(gtx, th)

			e.Frame(gtx.Ops)
		}
	}
}

// DrawGrid draws lines for the sudoku grid
func DrawGrid(gtx *layout.Context, th *material.Theme) {

}

// LayoutGrid creates a 2d nested List corresponding to the board
func LayoutGrid(gtx *layout.Context, th *material.Theme) {
	sections := make([]layout.FlexChild, 2)
	// Larger section consists of sudokuBoard elements in Lists
	sections[0] = layout.Flexed(0.9, func() {
		LayoutSudoku(gtx, th)
	})
	// Smaller section is for the program button and display
	sections[1] = layout.Flexed(0.1, func() {
		if !solving {
			LayoutStart(gtx, th)
		} else {
			LayoutEnd(gtx, th)
		}
	})
	layout.Flex{Axis: layout.Vertical}.Layout(gtx, sections...)
}

func LayoutStart(gtx *layout.Context, th *material.Theme) {
	for startButton.Clicked(gtx) {
		solving = true
		startTime = time.Now()
	}
	material.Clickable(gtx, startButton, func() {
		layout.Inset{
			Top:    unit.Px(20),
			Bottom: unit.Px(20),
			Left:   unit.Px(420),
			Right:  unit.Px(420),
		}.Layout(gtx, func() {
			msg := "Start"
			label := material.Caption(th, msg)
			label.Alignment = text.Middle
			label.Layout(gtx)
		})
	})
}

func LayoutEnd(gtx *layout.Context, th *material.Theme) {
	layout.Inset{
		Top:    unit.Px(20),
		Bottom: unit.Px(20),
		Left:   unit.Px(200),
		Right:  unit.Px(200),
	}.Layout(gtx, func() {
		var msg string
		if solved {
			msg = fmt.Sprintf("Took: %v", endTime)
		} else {
			msg = "Solving"
		}
		label := material.Caption(th, msg)
		label.Alignment = text.Middle
		label.Layout(gtx)
	})
}


func LayoutSudoku(gtx *layout.Context, th *material.Theme) {
	grid := &layout.List{
		Axis:        layout.Vertical,
		ScrollToEnd: false,
	}
	grid.Layout(gtx, 9, func(i int) {
		l := &layout.List{
			Axis:        layout.Horizontal,
			ScrollToEnd: false,
		}
		l.Layout(gtx, 9, func(j int) {
			layout.UniformInset(unit.Px(50)).Layout(gtx, func() {
				msg := strconv.Itoa(sudokuBoard[i][j])
				if msg == "0" {
					msg = " "
				}
				label := material.Caption(th, msg)
				label.Alignment = text.Middle
				label.Layout(gtx)
			})
		})
	})
}

// driver for solving algorithm
func main() {
	go SolveBoard(&sudokuBoard)
	CreateWindow()
}
