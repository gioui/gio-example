// SPDX-License-Identifier: Unlicense OR MIT

package main

import (
	"image"
	"log"
	"os"
	"time"

	"gioui.org/app"       // app contains Window handling.
	"gioui.org/io/key"    // key is used for keyboard events.
	"gioui.org/io/system" // system is used for system events (e.g. closing the window).
	"gioui.org/layout"    // layout is used for layouting widgets.
	"gioui.org/op"        // op is used for recording different operations.
	"gioui.org/unit"      // unit is used to define pixel-independent sizes
)

var (
	// cellSizePx is the cell size in pixels.
	cellSize = unit.Dp(5)
	// boardSize is the count of cells in a particular dimension.
	boardSize = image.Pt(50, 50)
)

func main() {
	// The ui loop is separated from the application window creation
	// such that it can be used for testing.
	ui := NewUI()

	windowWidth := cellSize.Scale(float32(boardSize.X + 2))
	windowHeight := cellSize.Scale(float32(boardSize.Y + 2))
	// This creates a new application window and starts the UI.
	go func() {
		w := app.NewWindow(
			app.Title("Game of Life"),
			app.Size(windowWidth, windowHeight),
		)
		if err := ui.Run(w); err != nil {
			log.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	}()

	// This starts Gio main.
	app.Main()
}

// UI holds all of the application state.
type UI struct {
	// Board handles all game-of-life logic.
	Board *Board
}

// NewUI creates a new UI using the Go Fonts.
func NewUI() *UI {
	// We start with a new random board.
	board := NewBoard(boardSize)
	board.Randomize()

	return &UI{
		Board: board,
	}
}

// Run handles window events and renders the application.
func (ui *UI) Run(w *app.Window) error {
	var ops op.Ops

	// Update the board 3 times per second.
	advanceBoard := time.NewTicker(time.Second / 3)
	defer advanceBoard.Stop()

	// listen for events happening on the window.
	for {
		select {
		case e := <-w.Events():
			// detect the type of the event.
			switch e := e.(type) {
			// this is sent when the application should re-render.
			case system.FrameEvent:
				// gtx is used to pass around rendering and event information.
				gtx := layout.NewContext(&ops, e)
				// render and handle UI.
				ui.Layout(gtx)
				// render and handle the operations from the UI.
				e.Frame(gtx.Ops)

			// handle a global key press.
			case key.Event:
				switch e.Name {
				// when we click escape, let's close the window.
				case key.NameEscape:
					return nil
				}

			// this is sent when the application is closed.
			case system.DestroyEvent:
				return e.Err
			}

		case <-advanceBoard.C:
			ui.Board.Advance()
			w.Invalidate()
		}
	}
}

// Layout displays the main program layout.
func (ui *UI) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Center.Layout(gtx,
		BoardStyle{
			CellSizePx: gtx.Px(cellSize),
			Board:      ui.Board,
		}.Layout,
	)
}
