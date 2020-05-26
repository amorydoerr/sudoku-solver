package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// CreateWindow opens a new gioui window
func CreateWindow() {
	gofont.Register()
	go func() {
		w := app.NewWindow(app.Size(unit.Dp(425), unit.Dp(475)), app.Title("Sudoku Solver"))
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
	style := material.ButtonLayoutStyle{
		Inset: layout.UniformInset(unit.Dp(10)),
	}
	style.Layout(gtx, startButton, func() {
		material.Clickable(gtx, startButton, func() {
			msg := "Start"
			label := material.Caption(th, msg)
			label.Alignment = text.Middle
			label.Layout(gtx)
		})
	})
}

func LayoutEnd(gtx *layout.Context, th *material.Theme) {
	layout.UniformInset(unit.Dp(10)).Layout(gtx, func() {
		var msg string
		if failed {
			msg = "No Solution"
		} else if solved {
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
			if !solving {
				LayoutInput(gtx, th, i, j)
			} else {
				LayoutValues(gtx, th, i, j)
			}
		})
	})
}

func LayoutValues(gtx *layout.Context, th *material.Theme, row, col int) {
	layout.Inset{
		Top:    unit.Dp(15),
		Bottom: unit.Dp(15),
		Left:   unit.Dp(20),
		Right:  unit.Dp(20),
	}.Layout(gtx, func() {
		msg := strconv.Itoa(sudokuBoard[row][col])
		if msg == "0" {
			msg = " "
		}
		label := material.Caption(th, msg)
		label.Alignment = text.Middle
		label.Layout(gtx)
	})
}

func LayoutInput(gtx *layout.Context, th *material.Theme, row, col int) {
	layout.Inset{
		Top:    unit.Dp(15),
		Bottom: unit.Dp(15),
		Left:   unit.Dp(20),
		Right:  unit.Dp(20),
	}.Layout(gtx, func() {
		field := new(widget.Editor)
		field.Alignment = text.Middle
		editor := material.Editor(th, "x")
		editor.Layout(gtx, field)

	})
}
