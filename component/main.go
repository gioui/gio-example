package main

import (
	"flag"
	"log"
	"os"

	"gioui.org/app"
	page "gioui.org/example/component/pages"
	"gioui.org/example/component/pages/about"
	"gioui.org/example/component/pages/appbar"
	"gioui.org/example/component/pages/discloser"
	"gioui.org/example/component/pages/menu"
	"gioui.org/example/component/pages/navdrawer"
	"gioui.org/example/component/pages/textfield"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget/material"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

func main() {
	flag.Parse()
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
	th := material.NewTheme(gofont.Collection())
	var ops op.Ops

	router := page.NewRouter()
	router.Register(0, appbar.New(&router))
	router.Register(1, navdrawer.New(&router))
	router.Register(2, textfield.New(&router))
	router.Register(3, menu.New(&router))
	router.Register(4, discloser.New(&router))
	router.Register(5, about.New(&router))

	for {
		select {
		case e := <-w.Events():
			switch e := e.(type) {
			case system.DestroyEvent:
				return e.Err
			case system.FrameEvent:
				gtx := layout.NewContext(&ops, e)
				router.Layout(gtx, th)
				e.Frame(gtx.Ops)
			}
		}
	}
}
