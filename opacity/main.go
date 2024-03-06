// SPDX-License-Identifier: Unlicense OR MIT

package main

import (
	"fmt"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/widget"
	"gioui.org/widget/material"
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

type opacityViewStyle struct {
	control *widget.Float
	slider  material.SliderStyle
	value   material.LabelStyle
	padding layout.Inset
}

func opacityView(th *material.Theme, state *widget.Float) opacityViewStyle {
	return opacityViewStyle{
		slider:  material.Slider(th, state),
		padding: layout.UniformInset(12),
		value:   material.Body1(th, fmt.Sprintf("%.2f", state.Value)),
	}
}

func (o opacityViewStyle) Layout(gtx C, w layout.Widget) D {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return o.padding.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return o.value.Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Width: o.padding.Left}.Layout),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return o.slider.Layout(gtx)
					}),
				)
			})
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			defer paint.PushOpacity(gtx.Ops, o.slider.Float.Value).Pop()
			return w(gtx)
		}),
	)
}

func loop(w *app.Window) error {
	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	var ops op.Ops
	var outer, inner widget.Float
	outer.Value = .75
	inner.Value = .5
	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			opacityView(th, &outer).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return opacityView(th, &inner).Layout(gtx,
					func(gtx layout.Context) layout.Dimensions {
						return material.Loader(th).Layout(gtx)
					})
			})
			e.Frame(gtx.Ops)
		}
	}
}
