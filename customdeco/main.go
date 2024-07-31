// SPDX-License-Identifier: Unlicense OR MIT

// The customdeco program demonstrates custom decorations
// in Gio.
package main

import (
	"fmt"
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
	"gioui.org/widget"
	"gioui.org/widget/material"

	"gioui.org/font/gofont"
)

func main() {
	go func() {
		w := new(app.Window)
		w.Option(app.Decorated(false))
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func loop(w *app.Window) error {
	var (
		b    widget.Clickable
		deco widget.Decorations
	)
	var (
		toggle    bool
		decorated bool
		title     string
	)
	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	var ops op.Ops
	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.ConfigEvent:
			decorated = e.Config.Decorated
			deco.Maximized = e.Config.Mode == app.Maximized
			title = e.Config.Title
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			for b.Clicked(gtx) {
				toggle = !toggle
				w.Option(app.Decorated(toggle))
			}
			cl := clip.Rect{Max: e.Size}.Push(gtx.Ops)
			paint.ColorOp{Color: color.NRGBA{A: 0xff, G: 0xff}}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(material.Button(th, &b, "Toggle decorations").Layout),
					layout.Rigid(material.Body1(th, fmt.Sprintf("Decorated: %v", decorated)).Layout),
				)
			})
			cl.Pop()
			if !decorated {
				w.Perform(deco.Update(gtx))
				material.Decorations(th, &deco, ^system.Action(0), title).Layout(gtx)
			}
			e.Frame(gtx.Ops)
		}
	}
}
