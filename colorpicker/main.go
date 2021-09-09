package main

import (
	"image"
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget/material"
	"gioui.org/x/colorpicker"
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

var white = color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}

func loop(w *app.Window) error {
	th := material.NewTheme(gofont.Collection())
	background := white
	current := color.NRGBA{R: 255, G: 128, B: 75, A: 255}
	picker := colorpicker.State{}
	picker.SetColor(current)
	muxState := colorpicker.NewMuxState(
		[]colorpicker.MuxOption{
			{
				Label: "current",
				Value: &current,
			},
			{
				Label: "background",
				Value: &th.Palette.Bg,
			},
			{
				Label: "foreground",
				Value: &th.Palette.Fg,
			},
		}...)
	background = *muxState.Color()
	var ops op.Ops
	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			log.Println("destroyed")
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			if muxState.Changed() {
				background = *muxState.Color()
				log.Printf("mux changed")
			}
			if picker.Changed() {
				current = picker.Color()
				background = *muxState.Color()
				log.Printf("picker changed")
			}
			layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return colorpicker.PickerStyle{Label: "Current", Theme: th, State: &picker}.Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return colorpicker.Mux(th, &muxState, "Display Right:").Layout(gtx)
						}),
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							size := gtx.Constraints.Max
							paint.FillShape(gtx.Ops, background, clip.Rect(image.Rectangle{Max: size}).Op())
							return D{Size: size}
						}),
					)
				}),
			)
			e.Frame(gtx.Ops)
		}
	}
}
