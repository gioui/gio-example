// SPDX-License-Identifier: Unlicense OR MIT

package main

// A simple Gio program. See https://gioui.org for more information.

import (
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"gioui.org/font/gofont"
	"gioui.org/x/haptic"
)

var buzzer *haptic.Buzzer

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

func loop(w *app.Window) error {
	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	btn := widget.Clickable{}
	buzzer = haptic.NewBuzzer(w)
	go func() {
		for err := range buzzer.Errors() {
			if err != nil {
				log.Printf("buzzer error: %v", err)
			}
		}
	}()
	var ops op.Ops
	for {
		switch e := w.NextEvent().(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			if btn.Clicked(gtx) {
				buzzer.Buzz()
			}
			layout.Center.Layout(gtx, material.Button(th, &btn, "buzz").Layout)
			e.Frame(gtx.Ops)
		default:
			ProcessPlatformEvent(e)
		}
	}
}
