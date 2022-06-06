// SPDX-License-Identifier: Unlicense OR MIT

package main

// A simple Gio program. See https://gioui.org for more information.

import (
	"fmt"
	"os"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"

	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"gioui.org/x/notify"

	"gioui.org/font/gofont"
)

func main() {
	go func() {
		w := app.NewWindow(
			app.Title("notify"),
			app.Size(unit.Dp(800), unit.Dp(600)))

		var ops op.Ops
		for event := range w.Events() {
			switch event := event.(type) {
			case system.DestroyEvent:
				os.Exit(0)
			case system.FrameEvent:
				event.Frame(frame(layout.NewContext(&ops, event)))
			}
		}
	}()
	app.Main()
}

type (
	// C quick alias for Context.
	C = layout.Context
	// D quick alias for Dimensions.
	D = layout.Dimensions
)

var (
	th       = material.NewTheme(gofont.Collection())
	notifier = func() notify.Notifier {
		n, err := notify.NewNotifier()
		if err != nil {
			panic(fmt.Errorf("init notification manager: %w", err))
		}
		return n
	}()
	editor    component.TextField
	notifyBtn widget.Clickable
)

// frame lays out the entire frame and returns the reusltant ops buffer.
func frame(gtx C) *op.Ops {
	if notifyBtn.Clicked() {
		msg := "This is a notification send from gio."
		if txt := editor.Text(); txt != "" {
			msg = txt
		}
		go notifier.CreateNotification("Hello Gio!", msg)
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
		)
	})
	return gtx.Ops
}
