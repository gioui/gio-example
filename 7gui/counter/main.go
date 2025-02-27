package main

import (
	"log"
	"os"
	"strconv"

	"gioui.org/app"         // app contains Window handling.
	"gioui.org/font/gofont" // gofont is used for loading the default font.
	"gioui.org/io/event"
	"gioui.org/io/key" // key is used for keyboard events.

	// system is used for system events (e.g. closing the window).
	"gioui.org/layout" // layout is used for layouting widgets.
	"gioui.org/op"     // op is used for recording different operations.
	"gioui.org/op/clip"
	"gioui.org/text"
	"gioui.org/unit"            // unit is used to define pixel-independent sizes
	"gioui.org/widget"          // widget contains state handling for widgets.
	"gioui.org/widget/material" // material contains material design widgets.
)

func main() {
	// The ui loop is separated from the application window creation
	// such that it can be used for testing.
	ui := NewUI()

	// This creates a new application window and starts the UI.
	go func() {
		w := new(app.Window)
		w.Option(
			app.Title("Counter"),
			app.Size(unit.Dp(240), unit.Dp(70)),
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

// defaultMargin is a margin applied in multiple places to give
// widgets room to breathe.
var defaultMargin = unit.Dp(10)

// UI holds all of the application state.
type UI struct {
	// Theme is used to hold the fonts used throughout the application.
	Theme *material.Theme

	// Counter displays and keeps the state of the counter.
	Counter Counter
}

// NewUI creates a new UI using the Go Fonts.
func NewUI() *UI {
	ui := &UI{}
	ui.Theme = material.NewTheme()
	ui.Theme.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	return ui
}

// Run handles window events and renders the application.
func (ui *UI) Run(w *app.Window) error {
	var ops op.Ops

	// listen for events happening on the window.
	for {
		// detect the type of the event.
		switch e := w.Event().(type) {
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
			return e.Err
		}
	}
}

// Layout displays the main program layout.
func (ui *UI) Layout(gtx layout.Context) layout.Dimensions {
	// inset is used to add padding around the window border.
	inset := layout.UniformInset(defaultMargin)
	return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return ui.Counter.Layout(ui.Theme, gtx)
	})
}

// Counter is a component that keeps track of it's state and
// displays itself as a label and a button.
type Counter struct {
	// Count is the current value.
	Count int

	// increase is used to track button clicks.
	increase widget.Clickable
}

// Layout lays out the counter and handles input.
func (counter *Counter) Layout(th *material.Theme, gtx layout.Context) layout.Dimensions {
	// first, update the count so that once it renders it has the correct value
	counter.UpdateCount(gtx)

	// Flex layout lays out widgets from left to right by default.
	return layout.Flex{}.Layout(gtx,
		// We use weight 1 for both text and count to make them the same size.
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			// We center align the text to the area available.
			return layout.Center.Layout(gtx,
				// Body1 is the default text size for reading.
				material.Body1(th, strconv.Itoa(counter.Count)).Layout)
		}),
		// We use an empty widget to add spacing between the text
		// and the button.
		layout.Rigid(layout.Spacer{Height: defaultMargin}.Layout),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			// display the button.
			return material.Button(th, &counter.increase, "Count").Layout(gtx)
		}),
	)
}

func (counter *Counter) UpdateCount(gtx layout.Context) {
	// For every click on the button increment the count.
	for counter.increase.Clicked(gtx) {
		counter.Count++
	}
}
