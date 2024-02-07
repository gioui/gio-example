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
	"gioui.org/font"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
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
	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	var (
		ops  op.Ops
		grid component.GridState
	)
	timingWindow := time.Second
	timings := []FrameTiming{}
	frameCounter := 0
	timingStart := time.Time{}
	for {
		switch e := w.NextEvent().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			gtx.Execute(op.InvalidateCmd{})
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

func layoutTable(th *material.Theme, gtx C, timings []FrameTiming, grid *component.GridState) D {
	// Configure width based on available space and a minimum size.
	minSize := gtx.Dp(unit.Dp(200))
	border := widget.Border{
		Color: color.NRGBA{A: 255},
		Width: unit.Dp(1),
	}

	inset := layout.UniformInset(unit.Dp(2))

	// Configure a label styled to be a heading.
	headingLabel := material.Body1(th, "")
	headingLabel.Font.Weight = font.Bold
	headingLabel.Alignment = text.Middle
	headingLabel.MaxLines = 1

	// Configure a label styled to be a data element.
	dataLabel := material.Body1(th, "")
	dataLabel.Font.Typeface = "Go Mono"
	dataLabel.MaxLines = 1
	dataLabel.Alignment = text.End

	// Measure the height of a heading row.
	orig := gtx.Constraints
	gtx.Constraints.Min = image.Point{}
	macro := op.Record(gtx.Ops)
	dims := inset.Layout(gtx, headingLabel.Layout)
	_ = macro.Stop()
	gtx.Constraints = orig

	return component.Table(th, grid).Layout(gtx, len(timings), 4,
		func(axis layout.Axis, index, constraint int) int {
			widthUnit := max(int(float32(constraint)/3), minSize)
			switch axis {
			case layout.Horizontal:
				switch index {
				case 0, 1:
					return int(widthUnit)
				case 2, 3:
					return int(widthUnit / 2)
				default:
					return 0
				}
			default:
				return dims.Size.Y
			}
		},
		func(gtx C, col int) D {
			return border.Layout(gtx, func(gtx C) D {
				return inset.Layout(gtx, func(gtx C) D {
					headingLabel.Text = headingText[col]
					return headingLabel.Layout(gtx)
				})
			})
		},
		func(gtx C, row, col int) D {
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
		},
	)
}
