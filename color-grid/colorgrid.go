// SPDX-License-Identifier: Unlicense OR MIT

package main

// A simple Gio program. See https://gioui.org for more information.

import (
	"image/color"
	"log"
	"os"

	"gioui.org/app"
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
		w := new(app.Window)
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
	th := material.NewTheme()
	var (
		ops  op.Ops
		grid component.GridState
	)
	sideLength := 1000
	cellSize := unit.Dp(10)
	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
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
