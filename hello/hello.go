// SPDX-License-Identifier: Unlicense OR MIT

package main

// A simple Gio program. See https://gioui.org for more information.

import (
	"fmt"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/io/event"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/op"
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

func loop(w *app.Window) error {
	tag := new(int)
	var ops op.Ops
	for {
		switch e := w.NextEvent().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			for {
				ev, ok := gtx.Source.Event(pointer.Filter{
					Target: tag,
					Kinds:  pointer.Release,
				})
				if !ok {
					break
				}
				switch ev := ev.(type) {
				case pointer.Event:
					if ev.Kind == pointer.Release {
						gtx.Execute(key.FocusCmd{Tag: tag})
						fmt.Println("triggered focus command")
					}
				}
				fmt.Printf("%#+v\n", ev)
			}
			for {
				ev, ok := gtx.Source.Event(key.Filter{
					Focus: tag,
				})
				if !ok {
					break
				}
				fmt.Printf("%#+v\n", ev)
			}
			event.Op(gtx.Ops, tag)
			e.Frame(gtx.Ops)
		}
	}
}
