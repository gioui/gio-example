// SPDX-License-Identifier: Unlicense OR MIT

package main

// A simple Gio program. See https://gioui.org for more information.

import (
	_ "embed"
	"fmt"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget"
	"gioui.org/widget/material"

	nsareg "eliasnaur.com/font/noto/sans/arabic/regular"
	"eliasnaur.com/font/roboto/robotoregular"
	"gioui.org/font/opentype"
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
	arabicFace, _ := opentype.Parse(nsareg.TTF)
	englishFace, _ := opentype.Parse(robotoregular.TTF)
	collection := []text.FontFace{}
	collection = append(collection, text.FontFace{Font: text.Font{Typeface: "Latin"}, Face: englishFace})
	collection = append(collection, text.FontFace{Font: text.Font{Typeface: "Arabic"}, Face: arabicFace})
	th := material.NewTheme(collection)
	var ed widget.Editor
	txt := "Hello أهلا my good friend صديقي الجيد bidirectional text نص ثنائي الاتجاه."
	ed.SetText(txt)
	ed.Focus()
	ed.Alignment = text.Middle
	var ops op.Ops
	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
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
