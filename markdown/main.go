// SPDX-License-Identifier: Unlicense OR MIT

package main

// A simple Gio program. See https://gioui.org for more information.
//
// This program showcases markdown rendering.
// The left pane contains a text editor for inputing raw text.
// The right pane renders the resulting markdown document using richtext.
//
// Richtext is fully interactive, links can be clicked, hovered, and longpressed.

import (
	"image"
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"gioui.org/x/markdown"
	"gioui.org/x/richtext"

	"gioui.org/font/gofont"
	"github.com/inkeliz/giohyperlink"
)

func main() {
	ui := UI{
		Window:   app.NewWindow(),
		Renderer: markdown.NewRenderer(),
		Shaper:   text.NewCache(gofont.Collection()),
		Theme:    NewTheme(gofont.Collection()),
		Resize:   component.Resize{Ratio: 0.5},
	}
	go func() {
		if err := ui.Loop(); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

type (
	C = layout.Context
	D = layout.Dimensions
)

// UI specifies the user interface.
type UI struct {
	// External systems.
	// Window provides access to the OS window.
	Window *app.Window
	// Theme contains semantic style data. Extends `material.Theme`.
	Theme *Theme
	// Shaper cache of registered fonts.
	Shaper *text.Cache
	// Renderer tranforms raw text containing markdown into richtext.
	Renderer *markdown.Renderer

	// Core state.
	// Editor retains raw text in an edit buffer.
	Editor widget.Editor
	// TextState retains rich text interactions: clicks, hovers and longpresses.
	TextState richtext.InteractiveText
	// Resize state retains the split between the editor and the rendered text.
	component.Resize
}

// Theme contains semantic style data.
type Theme struct {
	// Base theme to extend.
	Base *material.Theme
	// cache of processed markdown.
	cache []richtext.SpanStyle
}

// NewTheme instantiates a theme, extending material theme.
func NewTheme(font []text.FontFace) *Theme {
	return &Theme{
		Base: material.NewTheme(font),
	}
}

// Loop drives the UI until the window is destroyed.
func (ui UI) Loop() error {
	var ops op.Ops
	for {
		e := <-ui.Window.Events()
		giohyperlink.ListenEvents(e)
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			ui.Layout(gtx)
			e.Frame(gtx.Ops)
		}
	}
}

// Update processes events from the previous frame, updating state accordingly.
func (ui *UI) Update(gtx C) {
	for o, events := ui.TextState.Events(); o != nil; o, events = ui.TextState.Events() {
		for _, e := range events {
			switch e.Type {
			case richtext.Click:
				if url := o.Get(markdown.MetadataURL); url != "" {
					if err := giohyperlink.Open(url); err != nil {
						// TODO(jfm): display UI element explaining the error to the user.
						log.Printf("error: opening hyperlink: %v", err)
					}
				}
			case richtext.Hover:
			case richtext.LongPress:
				log.Println("longpress")
				ui.Window.Option(app.Title(o.Get(markdown.MetadataURL)))
			}
		}
	}
	for _, event := range ui.Editor.Events() {
		if _, ok := event.(widget.ChangeEvent); ok {
			var err error
			ui.Theme.cache, err = ui.Renderer.Render(ui.Theme.Base, []byte(ui.Editor.Text()))
			if err != nil {
				// TODO(jfm): display UI element explaining the error to the user.
				log.Printf("error: rendering markdown: %v", err)
			}
		}
	}
}

// Layout renders the current frame.
func (ui *UI) Layout(gtx C) D {
	ui.Update(gtx)
	return ui.Resize.Layout(gtx,
		func(gtx C) D {
			return layout.UniformInset(unit.Dp(4)).Layout(gtx, func(gtx C) D {
				return material.Editor(ui.Theme.Base, &ui.Editor, "markdown").Layout(gtx)
			})
		},
		func(gtx C) D {
			return layout.UniformInset(unit.Dp(4)).Layout(gtx, func(gtx C) D {
				return richtext.Text(&ui.TextState, ui.Shaper, ui.Theme.cache...).Layout(gtx)
			})
		},
		func(gtx C) D {
			rect := image.Rectangle{
				Max: image.Point{
					X: (gtx.Dp(unit.Dp(4))),
					Y: (gtx.Constraints.Max.Y),
				},
			}
			paint.FillShape(gtx.Ops, color.NRGBA{A: 200}, clip.Rect(rect).Op())
			return D{Size: rect.Max}
		},
	)
}
