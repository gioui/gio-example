// SPDX-License-Identifier: Unlicense OR MIT

package main

import (
	"image"
	"log"
	"os"
	"time"

	"gioui.org/app" // app contains Window handling.
	"gioui.org/io/event"
	"gioui.org/io/key" // key is used for keyboard events.

	// system is used for system events (e.g. closing the window).
	"gioui.org/layout" // layout is used for layouting widgets.
	"gioui.org/op"     // op is used for recording different operations.
	"gioui.org/op/clip"
	"gioui.org/unit" // unit is used to define pixel-independent sizes
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

	windowWidth := cellSize * (unit.Dp(boardSize.X + 2))
	windowHeight := cellSize * (unit.Dp(boardSize.Y + 2))
	// This creates a new application window and starts the UI.
	go func() {
		w := new(app.Window)
		w.Option(
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

	events := make(chan event.Event)
	acks := make(chan struct{})

	go func() {
		for {
			ev := w.Event()
			events <- ev
			<-acks
			if _, ok := ev.(app.DestroyEvent); ok {
				return
			}
		}
	}()

	// listen for events happening on the window.
	for {
		select {
		case e := <-events:
			// detect the type of the event.
			switch e := e.(type) {
			// this is sent when the application should re-render.
			case app.FrameEvent:
				// gtx is used to pass around rendering and event information.
				gtx := app.NewContext(&ops, e)
				// register a global key listener for the escape key wrapping our entire UI.
				area := clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops)
				event.Op(gtx.Ops, w)

				// check for presses of the escape key and close the window if we find them.
				for {
					event, ok := gtx.Event(key.Filter{
						Name: key.NameEscape,
					})
					if !ok {
						break
					}
					switch event := event.(type) {
					case key.Event:
						if event.Name == key.NameEscape {
							return nil
						}
					}
				}
				// render and handle UI.
				ui.Layout(gtx)
				area.Pop()
				// render and handle the operations from the UI.
				e.Frame(gtx.Ops)

			// this is sent when the application is closed.
			case app.DestroyEvent:
				acks <- struct{}{}
				return e.Err
			}
			acks <- struct{}{}

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
			CellSizePx: gtx.Dp(cellSize),
			Board:      ui.Board,
		}.Layout,
	)
}
