// SPDX-License-Identifier: Unlicense OR MIT

package main

// A simple Gio program. See https://gioui.org for more information.

import (
	_ "embed"
	"fmt"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func main() {
	go func() {
		w := new(app.Window)
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
	th := material.NewTheme()
	var ed widget.Editor
	txt := "Hello أهلا my good friend صديقي الجيد bidirectional text نص ثنائي الاتجاه."
	ed.SetText(txt)
	init := false
	ed.Alignment = text.Middle
	var ops op.Ops
	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			if !init {
				init = true
				gtx.Execute(key.FocusCmd{Tag: &ed})
			}
			layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					med := material.Editor(th, &ed, "")
					med.TextSize = material.H3(th, "").TextSize
					return med.Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					start, end := ed.Selection()
					return material.Body1(th, fmt.Sprintf("Selection start %d, end %d", start, end)).Layout(gtx)
				}),
			)
			e.Frame(gtx.Ops)
		}
	}
}
