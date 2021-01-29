// SPDX-License-Identifier: Unlicense OR MIT

package main

// A simple Gio program. See https://gioui.org for more information.

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"reflect"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"gioui.org/font/gofont"

	"gioui.org/example/livedit/giopkgs"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
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

const startingText = `package live

import (
    "gioui.org/layout"
    "gioui.org/widget/material"
    "gioui.org/op"
    "gioui.org/f32"

    "math"
)

var rotation float32

func Layout(gtx layout.Context, theme *material.Theme) layout.Dimensions {
    rotation += math.Pi/120
    op.Affine(f32.Affine2D{}.Rotate(f32.Point{}, rotation)).Add(gtx.Ops)
    op.InvalidateOp{}.Add(gtx.Ops)
    return material.H1(theme, "Hello!").Layout(gtx)
}
`

func loop(w *app.Window) error {
	th := material.NewTheme(gofont.Collection())
	var editor widget.Editor
	editor.SetText(startingText)
	var ops op.Ops
	var yaegi *interp.Interpreter

	first := true
	maroon := color.NRGBA{R: 127, G: 0, B: 0, A: 255}
	var (
		custom func(C, *material.Theme) D
		err    error
	)

	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			if len(editor.Events()) > 0 || first {
				yaegi = interp.New(interp.Options{})
				yaegi.Use(stdlib.Symbols)
				yaegi.Use(interp.Symbols)
				yaegi.Use(giopkgs.Symbols)
				func() {
					_, err = yaegi.Eval(editor.Text())
					if err != nil {
						log.Println(err)
						return
					}
					var result reflect.Value
					result, err = yaegi.Eval("live.Layout")
					if err != nil {
						log.Println(err)
						return
					}
					var (
						ok        bool
						newCustom func(C, *material.Theme) D
					)
					newCustom, ok = result.Interface().(func(layout.Context, *material.Theme) layout.Dimensions)
					if !ok {
						err = fmt.Errorf("returned data is not a widget, is %s", result.Type())
						log.Println(err)
						return
					}
					custom = newCustom
				}()
			}

			layout.Flex{}.Layout(gtx,
				layout.Flexed(.5, func(gtx C) D {
					return layout.UniformInset(unit.Dp(8)).Layout(gtx, material.Editor(th, &editor, "write layout code here").Layout)
				}),
				layout.Flexed(.5, func(gtx C) D {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							if err != nil {
								msg := material.Body1(th, err.Error())
								msg.Color = maroon
								return msg.Layout(gtx)
							}
							return D{}
						}),
						layout.Flexed(1.0, func(gtx C) (dims D) {
							defer func() {
								if err := recover(); err != nil {
									msg := material.Body1(th, "panic: "+err.(error).Error())
									msg.Color = maroon
									dims = msg.Layout(gtx)
								}
							}()
							if custom == nil {
								msg := material.Body1(th, "nil")
								msg.Color = maroon
								return msg.Layout(gtx)
							}
							return custom(gtx, th)
						}),
					)
				}),
			)

			e.Frame(gtx.Ops)
			first = false
		}
	}
}
