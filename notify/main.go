// SPDX-License-Identifier: Unlicense OR MIT

package main

// A simple Gio program. See https://gioui.org for more information.

import (
	"fmt"
	"os"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"

	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"gioui.org/x/notify"

	"gioui.org/font/gofont"
)

type (
	// C quick alias for Context.
	C = layout.Context
	// D quick alias for Dimensions.
	D = layout.Dimensions
)

func main() {
	go func() {
		th := material.NewTheme()
		th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
		n, err := notify.NewNotifier()
		if err != nil {
			panic(fmt.Errorf("init notification manager: %w", err))
		}
		_, ongoingSupported := n.(notify.OngoingNotifier)
		notifier := n
		var editor component.TextField
		var notifyBtn widget.Clickable
		var setOngoing widget.Bool
		w := new(app.Window)
		w.Option(
			app.Title("notify"),
			app.Size(unit.Dp(800), unit.Dp(600)))

		var ops op.Ops
		for {
			switch event := w.Event().(type) {
			case app.DestroyEvent:
				os.Exit(0)
			case app.FrameEvent:
				gtx := app.NewContext(&ops, event)
				if notifyBtn.Clicked(gtx) {
					msg := "This is a notification send from gio."
					if txt := editor.Text(); txt != "" {
						msg = txt
					}
					if ongoingSupported && setOngoing.Value {
						go notifier.(notify.OngoingNotifier).CreateOngoingNotification("Hello Gio!", msg)
					} else {
						go notifier.CreateNotification("Hello Gio!", msg)
					}
				}
				layout.Center.Layout(gtx, func(gtx C) D {
					gtx.Constraints.Max.X = gtx.Dp(unit.Dp(300))
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return editor.Layout(gtx, th, "enter a notification message")
						}),
						layout.Rigid(func(gtx C) D {
							return layout.Spacer{Height: unit.Dp(10)}.Layout(gtx)
						}),
						layout.Rigid(func(gtx C) D {
							return material.Button(th, &notifyBtn, "notify").Layout(gtx)
						}),
						layout.Rigid(func(gtx C) D {
							return layout.Spacer{Height: unit.Dp(10)}.Layout(gtx)
						}),
						layout.Rigid(func(gtx C) D {
							if !ongoingSupported {
								return D{}
							}
							return material.CheckBox(th, &setOngoing, "ongoing").Layout(gtx)
						}),
					)
				})
				event.Frame(gtx.Ops)
			}
		}
	}()
	app.Main()
}
