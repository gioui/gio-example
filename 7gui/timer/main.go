package main

import (
	"log"
	"os"
	"time"

	"gioui.org/app"             // app contains Window handling.
	"gioui.org/font/gofont"     // gofont is used for loading the default font.
	"gioui.org/io/key"          // key is used for keyboard events.
	"gioui.org/io/system"       // system is used for system events (e.g. closing the window).
	"gioui.org/layout"          // layout is used for layouting widgets.
	"gioui.org/op"              // op is used for recording different operations.
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
		w := app.NewWindow(
			app.Title("Timer"),
			app.Size(unit.Dp(360), unit.Dp(360)),
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

	Timer *Timer

	duration widget.Float
	reset    widget.Clickable
}

// NewUI creates a new UI using the Go Fonts.
func NewUI() *UI {
	ui := &UI{}
	ui.Theme = material.NewTheme(gofont.Collection())

	// start with reasonable defaults.
	ui.Timer = NewTimer(5 * time.Second)
	ui.duration.Value = 5

	return ui
}

// Run handles window events and renders the application.
func (ui *UI) Run(w *app.Window) error {

	// start the timer goroutine and ensure it's closed
	// when the application closes.
	closeTimer := ui.Timer.Start()
	defer closeTimer()

	var ops op.Ops
	for {
		select {
		// when the timer is updated we should update the screen.
		case <-ui.Timer.Updated:
			w.Invalidate()

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
		}
	}

	return nil
}

// Layout displays the main program layout.
func (ui *UI) Layout(gtx layout.Context) layout.Dimensions {
	th := ui.Theme

	// check whether the reset button was clicked.
	if ui.reset.Clicked() {
		ui.Timer.Reset()
	}
	// check whether the slider value has changed.
	if ui.duration.Changed() {
		ui.Timer.SetDuration(secondsToDuration(float64(ui.duration.Value)))
	}

	// get the latest information about the timer.
	info := ui.Timer.Info()
	progress := float32(0)
	if info.Duration == 0 {
		progress = 1
	} else {
		progress = float32(info.Progress.Seconds() / info.Duration.Seconds())
	}

	// inset is used to add padding around the window border.
	inset := layout.UniformInset(defaultMargin)
	return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(material.Body1(th, "Elapsed Time").Layout),
			layout.Rigid(material.ProgressBar(th, progress).Layout),
			layout.Rigid(material.Body1(th, info.ProgressString()).Layout),

			layout.Rigid(layout.Spacer{Height: th.TextSize}.Layout),
			layout.Rigid(material.Body1(th, "Duration").Layout),
			layout.Rigid(material.Slider(th, &ui.duration, 0, 15).Layout),

			layout.Rigid(layout.Spacer{Height: th.TextSize}.Layout),
			layout.Rigid(material.Button(th, &ui.reset, "Reset").Layout),
		)
	})
}

func secondsToDuration(s float64) time.Duration {
	return time.Duration(s * float64(time.Second))
}
