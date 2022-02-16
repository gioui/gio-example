// SPDX-License-Identifier: Unlicense OR MIT

package main

// A simple Gio program. See https://gioui.org for more information.

import (
	"image"
	"image/color"
	"log"
	"os"
	"strconv"
	"time"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
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

func calcWidths(gtx C, widths []unit.Value, quantity int, size unit.Value) []unit.Value {
	widths = widths[:0]
	for i := 0; i < quantity; i++ {
		widths = append(widths, size)
	}
	return widths
}

type FrameTiming struct {
	Start, End      time.Time
	FrameCount      int
	FramesPerSecond float64
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func loop(w *app.Window) error {
	th := material.NewTheme(gofont.Collection())
	var (
		ops  op.Ops
		grid widget.Grid
	)
	timingWindow := time.Second
	timings := []FrameTiming{}
	frameCounter := 0
	timingStart := time.Time{}
	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			op.InvalidateOp{}.Add(gtx.Ops)
			if timingStart == (time.Time{}) {
				timingStart = gtx.Now
			}
			if interval := gtx.Now.Sub(timingStart); interval >= timingWindow {
				timings = append(timings, FrameTiming{
					Start:           timingStart,
					End:             gtx.Now,
					FrameCount:      frameCounter,
					FramesPerSecond: float64(frameCounter) / interval.Seconds(),
				})
				frameCounter = 0
				timingStart = gtx.Now
			}
			layoutTable(th, gtx, timings, &grid)
			e.Frame(gtx.Ops)
			frameCounter++
		}
	}
}

var headingText = []string{"Start", "End", "Frames", "FPS"}

func layoutTable(th *material.Theme, gtx C, timings []FrameTiming, grid *widget.Grid) D {
	// Configure width based on available space and a minimum size.
	widthUnit := float32(max(gtx.Constraints.Max.X/3, gtx.Px(unit.Dp(200))))
	widths := []unit.Value{
		unit.Px(widthUnit),
		unit.Px(widthUnit),
		unit.Px(widthUnit * .5),
		unit.Px(widthUnit * .5),
	}
	border := widget.Border{
		Color: color.NRGBA{A: 255},
		Width: unit.Px(1),
	}

	inset := layout.UniformInset(unit.Dp(2))

	// Configure a label styled to be a heading.
	headingLabel := material.Body1(th, "")
	headingLabel.Font.Weight = text.Bold
	headingLabel.Alignment = text.Middle
	headingLabel.MaxLines = 1

	// Configure a label styled to be a data element.
	dataLabel := material.Body1(th, "")
	dataLabel.Font.Variant = "Mono"
	dataLabel.MaxLines = 1
	dataLabel.Alignment = text.End

	// Measure the height of a heading row.
	orig := gtx.Constraints
	gtx.Constraints.Min = image.Point{}
	macro := op.Record(gtx.Ops)
	dims := inset.Layout(gtx, headingLabel.Layout)
	_ = macro.Stop()
	cellHeight := unit.Px(float32(dims.Size.Y))
	gtx.Constraints = orig

	return material.Table(th, grid).Layout(gtx, len(timings), cellHeight, widths, func(gtx C, row, col int) D {
		return inset.Layout(gtx, func(gtx C) D {
			timing := timings[row]
			switch col {
			case 0:
				dataLabel.Text = timing.Start.Format("15:04:05.000000")
			case 1:
				dataLabel.Text = timing.End.Format("15:04:05.000000")
			case 2:
				dataLabel.Text = strconv.Itoa(timing.FrameCount)
			case 3:
				dataLabel.Text = strconv.FormatFloat(timing.FramesPerSecond, 'f', 2, 64)
			}
			return dataLabel.Layout(gtx)
		})
	}, func(gtx C, col int) D {
		return border.Layout(gtx, func(gtx C) D {
			return inset.Layout(gtx, func(gtx C) D {
				headingLabel.Text = headingText[col]
				return headingLabel.Layout(gtx)
			})
		})
	})
}
