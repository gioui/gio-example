// SPDX-License-Identifier: Unlicense OR MIT

package main

// A simple Gio program. See https://gioui.org for more information.

import (
	"log"
	"math"
	"os"
	"strconv"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget"
	"gioui.org/widget/material"

	colEmoji "eliasnaur.com/font/noto/emoji/color"
	"gioui.org/font/gofont"
	"gioui.org/font/opentype"
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

func loop(w *app.Window) error {
	// Load default font collection.
	collection := gofont.Collection()
	// Load a color emoji font.
	faces, err := opentype.ParseCollection(colEmoji.TTF)
	if err != nil {
		panic(err)
	}
	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(append(collection, faces...)))
	var ops op.Ops
	var sel widget.Selectable
	message := "🥳🧁🍰🎁🎂🎈🎺🎉🎊\n📧〽️🧿🌶️🔋\n😂❤️😍🤣😊\n🥺🙏💕😭😘\n👍😅👏"
	var customTruncator widget.Bool
	var maxLines widget.Float
	maxLines.Value = 0

	const (
		minLinesRange = 1
		maxLinesRange = 5
	)
	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			inset := layout.UniformInset(5)
			layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					l := material.H4(th, message)
					if customTruncator.Value {
						l.Truncator = "cont..."
					} else {
						l.Truncator = ""
					}
					l.MaxLines = minLinesRange + int(math.Round(float64(maxLines.Value)*(maxLinesRange-minLinesRange)))
					l.State = &sel
					return inset.Layout(gtx, l.Layout)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{}.Layout(gtx,
							layout.Rigid(material.Switch(th, &customTruncator, "Use Custom Truncator").Layout),
							layout.Rigid(layout.Spacer{Width: 5}.Layout),
							layout.Rigid(material.Body1(th, "Use Custom Truncator").Layout),
						)
					})
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return inset.Layout(gtx, material.Body1(th, "Max Lines:").Layout)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
							layout.Rigid(material.Body2(th, strconv.Itoa(minLinesRange)).Layout),
							layout.Rigid(layout.Spacer{Width: 5}.Layout),
							layout.Flexed(1, material.Slider(th, &maxLines).Layout),
							layout.Rigid(layout.Spacer{Width: 5}.Layout),
							layout.Rigid(material.Body2(th, strconv.Itoa(maxLinesRange)).Layout),
						)
					})
				}),
			)
			e.Frame(gtx.Ops)
		}
	}
}
