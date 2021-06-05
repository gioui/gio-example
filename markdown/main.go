// SPDX-License-Identifier: Unlicense OR MIT

package main

// A simple Gio program. See https://gioui.org for more information.

import (
	"image"
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/gesture"
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
	go func() {
		w := app.NewWindow()
		if err := loop(w); err != nil {
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

func loop(w *app.Window) error {
	fontCollection := gofont.Collection()
	shaper := text.NewCache(fontCollection)
	th := material.NewTheme(fontCollection)
	renderer := markdown.NewRenderer()
	var ops op.Ops

	var ed widget.Editor
	var rs component.Resize
	rs.Ratio = .5
	var textState richtext.InteractiveText
	var rendered []richtext.SpanStyle
	inset := layout.UniformInset(unit.Dp(4))
	for {
		e := <-w.Events()
		giohyperlink.ListenEvents(e)
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			if to := textState.LongPressed(); to != nil {
				w.Option(app.Title(to.Get(markdown.MetadataURL)))
			}
			for o, events := textState.Events(gtx); o != nil; o, events = textState.Events(gtx) {
				for _, e := range events {
					switch e.Type {
					case gesture.TypeClick:
						if url := o.Get(markdown.MetadataURL); url != "" {
							giohyperlink.Open(url)
						}
					}
				}
			}

			for _, edEvent := range ed.Events() {
				if _, ok := edEvent.(widget.ChangeEvent); ok {
					rendered, _ = renderer.Render(th, []byte(ed.Text()))
				}
			}

			rs.Layout(gtx,
				func(gtx C) D { return inset.Layout(gtx, material.Editor(th, &ed, "markdown").Layout) },
				func(gtx C) D {
					return inset.Layout(gtx, func(gtx C) D {
						return richtext.Text(&textState, rendered...).Layout(gtx, shaper)
					})
				},
				func(gtx C) D {
					rect := image.Rectangle{
						Max: image.Point{
							X: (gtx.Px(unit.Dp(4))),
							Y: (gtx.Constraints.Max.Y),
						},
					}
					paint.FillShape(gtx.Ops, color.NRGBA{A: 200}, clip.Rect(rect).Op())
					return D{Size: rect.Max}
				},
			)
			e.Frame(gtx.Ops)
		}
	}
}
