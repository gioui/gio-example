// SPDX-License-Identifier: Unlicense OR MIT

package main

// A simple Gio program. See https://gioui.org for more information.

import (
	//	"image/color"
	"log"
	"time"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"

	//	"gioui.org/text"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"git.sr.ht/~whereswaldon/niotify"

	"gioui.org/font/gofont"
)

func main() {
	go func() {
		w := app.NewWindow()
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
	}()
	app.Main()
}

func loop(w *app.Window) error {
	th := material.NewTheme(gofont.Collection())
	var ops op.Ops
	first := true
	notificationRequests := make(chan struct{})
	var button widget.Clickable
	var err error
	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			if button.Clicked() {
				notificationRequests <- struct{}{}
			}
			gtx := layout.NewContext(&ops, e)
			layout.Inset{
    			Top: e.Insets.Top,
    			Bottom: e.Insets.Bottom,
    			Left: e.Insets.Left,
    			Right: e.Insets.Right,
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {

			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
                              layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                                           text := "notification errors will appear here"
                                           if err != nil {
                                               text = err.Error()
                                           }
                                   return material.Body1(th, text).Layout(gtx)
                              }),
                              layout.Flexed(1,func(gtx layout.Context) layout.Dimensions {
                                    return material.Button(th, &button, "Send Notification").Layout(gtx)
                              }),
                          )
			})
			e.Frame(gtx.Ops)
			if first {
				first = false
				go func() {
					mgr, e := niotify.NewManager()
					if e != nil {
						log.Printf("manager creation failed: %v", e)
						err = e
					}
					for _ = range notificationRequests {
    						log.Println("trying to send notification")
						notif, e := mgr.CreateNotification("hello!", "IS GIO OUT THERE?")
						if e != nil {
							log.Printf("notification send failed: %v", e)
							err = e
							continue
						}
						go func() {
							time.Sleep(time.Second * 10)
							if err = notif.Cancel(); err != nil {
								log.Printf("failed cancelling: %v", err)
							}
						}()
					}
				}()
			}
		}
	}
}
