package main

import (
	"image/color"
	"log"
	"os"
	"strconv"

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
			app.Title("Temperature Converter"),
			app.Size(unit.Dp(360), unit.Dp(47)),
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

	// Converter displays and modifies the state.
	Converter Converter
}

// NewUI creates a new UI using the Go Fonts.
func NewUI() *UI {
	ui := &UI{}
	ui.Theme = material.NewTheme(gofont.Collection())

	ui.Converter.Init()
	return ui
}

// Run handles window events and renders the application.
func (ui *UI) Run(w *app.Window) error {
	var ops op.Ops

	// listen for events happening on the window.
	for e := range w.Events() {
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

	return nil
}

// Layout displays the main program layout.
func (ui *UI) Layout(gtx layout.Context) layout.Dimensions {
	// inset is used to add padding around the window border.
	inset := layout.UniformInset(defaultMargin)
	return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return ui.Converter.Layout(ui.Theme, gtx)
	})
}

// Converter is a component that keeps track of it's state and
// displays itself as two editors.
type Converter struct {
	Celsius    Field
	Fahrenheit Field
}

// Init is used to set the inital state.
func (conv *Converter) Init() {
	conv.Celsius.SingleLine = true
	conv.Fahrenheit.SingleLine = true
}

// Layout lays out the editors.
func (conv *Converter) Layout(th *material.Theme, gtx layout.Context) layout.Dimensions {
	// We use an empty widget to add spacing between widgets.
	spacer := layout.Rigid(layout.Spacer{Width: defaultMargin}.Layout)

	// check whether the celsius value has changed.
	if conv.Celsius.Changed() {
		// try to convert the value to an integer
		newValue, err := strconv.Atoi(conv.Celsius.Text())
		// update whether the editor is displaying a valid value
		conv.Celsius.Invalid = err != nil
		if !conv.Celsius.Invalid {
			// update the other editor when it's valid
			conv.Fahrenheit.Invalid = false
			conv.Fahrenheit.SetText(strconv.Itoa(newValue*9/5 + 32))
		}
	}

	// check whether the fahrenheit value has changed.
	if conv.Fahrenheit.Changed() {
		newValue, err := strconv.Atoi(conv.Fahrenheit.Text())
		conv.Fahrenheit.Invalid = err != nil
		if !conv.Fahrenheit.Invalid {
			conv.Celsius.Invalid = false
			conv.Celsius.SetText(strconv.Itoa((newValue - 32) * 5 / 9))
		}
	}

	// TODO: use proper baseline alignment.
	return layout.Flex{}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return conv.Celsius.Layout(th, gtx)
		}),
		spacer,
		layout.Rigid(material.Body1(th, "Celsius").Layout),
		spacer,
		layout.Rigid(material.Body1(th, "=").Layout),
		spacer,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return conv.Fahrenheit.Layout(th, gtx)
		}),
		spacer,
		layout.Rigid(material.Body1(th, "Fahrenheit").Layout),
	)
}

// Field implements an editor that allows updating the state and detect
// changes to the field from other sources.
type Field struct {
	widget.Editor
	Invalid bool

	old string
}

// Changed checks once whether the editor context has changed.
func (ed *Field) Changed() bool {
	newText := ed.Editor.Text()
	changed := newText != ed.old
	ed.old = newText
	return changed
}

// SetText sets editor content without marking the editor changed.
func (ed *Field) SetText(s string) {
	ed.old = s
	ed.Editor.SetText(s)
}

// Layout handles the editor with the appropriate color and border.
func (ed *Field) Layout(th *material.Theme, gtx layout.Context) layout.Dimensions {
	// Determine colors based on the state of the editor.
	borderWidth := float32(0.5)
	borderColor := color.NRGBA{A: 107}
	switch {
	case ed.Editor.Focused():
		borderColor = th.Palette.ContrastBg
		borderWidth = 2
	case ed.Invalid:
		borderColor = color.NRGBA{R: 200, A: 0xFF}
	}

	// draw an editor with a border.
	return widget.Border{
		Color:        borderColor,
		CornerRadius: unit.Dp(4),
		Width:        unit.Dp(borderWidth),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(unit.Dp(4)).Layout(gtx,
			material.Editor(th, &ed.Editor, "").Layout)
	})
}
