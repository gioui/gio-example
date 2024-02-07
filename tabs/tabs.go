// SPDX-License-Identifier: Unlicense OR MIT

package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"os"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"gioui.org/font/gofont"
)

func main() {
	go func() {
		defer os.Exit(0)
		w := app.NewWindow()
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
	}()
	app.Main()
}

func loop(w *app.Window) error {
	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	var ops op.Ops
	for {
		switch e := w.NextEvent().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			drawTabs(gtx, th)
			e.Frame(gtx.Ops)
		}
	}
}

var (
	tabs   Tabs
	slider Slider
)

type Tabs struct {
	list     layout.List
	tabs     []Tab
	selected int
}

type Tab struct {
	btn   widget.Clickable
	Title string
}

func init() {
	for i := 1; i <= 100; i++ {
		tabs.tabs = append(tabs.tabs,
			Tab{Title: fmt.Sprintf("Tab %d", i)},
		)
	}
}

type (
	C = layout.Context
	D = layout.Dimensions
)

func drawTabs(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return tabs.list.Layout(gtx, len(tabs.tabs), func(gtx C, tabIdx int) D {
				t := &tabs.tabs[tabIdx]
				if t.btn.Clicked(gtx) {
					if tabs.selected < tabIdx {
						slider.PushLeft()
					} else if tabs.selected > tabIdx {
						slider.PushRight()
					}
					tabs.selected = tabIdx
				}
				var tabWidth int
				return layout.Stack{Alignment: layout.S}.Layout(gtx,
					layout.Stacked(func(gtx C) D {
						dims := material.Clickable(gtx, &t.btn, func(gtx C) D {
							return layout.UniformInset(unit.Dp(12)).Layout(gtx,
								material.H6(th, t.Title).Layout,
							)
						})
						tabWidth = dims.Size.X
						return dims
					}),
					layout.Stacked(func(gtx C) D {
						if tabs.selected != tabIdx {
							return layout.Dimensions{}
						}
						tabHeight := gtx.Dp(unit.Dp(4))
						tabRect := image.Rect(0, 0, tabWidth, tabHeight)
						paint.FillShape(gtx.Ops, th.Palette.ContrastBg, clip.Rect(tabRect).Op())
						return layout.Dimensions{
							Size: image.Point{X: tabWidth, Y: tabHeight},
						}
					}),
				)
			})
		}),
		layout.Flexed(1, func(gtx C) D {
			return slider.Layout(gtx, func(gtx C) D {
				fill(gtx, dynamicColor(tabs.selected), dynamicColor(tabs.selected+1))
				return layout.Center.Layout(gtx,
					material.H1(th, fmt.Sprintf("Tab content #%d", tabs.selected+1)).Layout,
				)
			})
		}),
	)
}

func fill(gtx layout.Context, col1, col2 color.NRGBA) {
	dr := image.Rectangle{Max: gtx.Constraints.Min}
	paint.FillShape(gtx.Ops,
		color.NRGBA{R: 0, G: 0, B: 0, A: 0xFF},
		clip.Rect(dr).Op(),
	)

	col2.R = byte(float32(col2.R))
	col2.G = byte(float32(col2.G))
	col2.B = byte(float32(col2.B))
	paint.LinearGradientOp{
		Stop1:  f32.Pt(float32(dr.Min.X), 0),
		Stop2:  f32.Pt(float32(dr.Max.X), 0),
		Color1: col1,
		Color2: col2,
	}.Add(gtx.Ops)
	defer clip.Rect(dr).Push(gtx.Ops).Pop()
	paint.PaintOp{}.Add(gtx.Ops)
}

func dynamicColor(i int) color.NRGBA {
	sn, cs := math.Sincos(float64(i) * math.Phi)
	return color.NRGBA{
		R: 0xA0 + byte(0x30*sn),
		G: 0xA0 + byte(0x30*cs),
		B: 0xD0,
		A: 0xFF,
	}
}
