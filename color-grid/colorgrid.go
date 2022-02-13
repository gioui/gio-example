// SPDX-License-Identifier: Unlicense OR MIT

package main

// A simple Gio program. See https://gioui.org for more information.

import (
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
	"gioui.org/unit"
	"gioui.org/widget/material"
	"gioui.org/x/component"
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
	th := material.NewTheme(gofont.Collection())
	var (
		ops  op.Ops
		grid component.GridState
	)
	sideLength := 1000
	cellSize := unit.Dp(10)
	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			component.Grid(th, &grid).Layout(gtx, sideLength, sideLength,
				func(axis layout.Axis, index, constraint int) int {
					return gtx.Dp(cellSize)
				},
				func(gtx C, row, col int) D {
					c := color.NRGBA{R: uint8(3 * row), G: uint8(5 * col), B: uint8(row * col), A: 255}
					paint.FillShape(gtx.Ops, c, clip.Rect{Max: gtx.Constraints.Max}.Op())
					return D{Size: gtx.Constraints.Max}
				})
			e.Frame(gtx.Ops)
		}
	}
}
