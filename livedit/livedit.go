// SPDX-License-Identifier: Unlicense OR MIT

package main

// A simple Gio program. See https://gioui.org for more information.

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"reflect"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"gioui.org/font/gofont"

	"gioui.org/example/livedit/giopkgs"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"

	"golang.org/x/image/draw"
	_ "image/png"
)

func main() {
	go func() {
		w := app.NewWindow(app.Title("Gio Playground"))
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

func loadImage(name string) image.Image {
	file, err := os.Open(name)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		log.Println(err)
		return nil
	}
	return img
}

func imageOpOrEmpty(src image.Image) paint.ImageOp {
	if src == nil {
		return paint.ImageOp{}
	}
	return paint.NewImageOp(src)
}

var (
	yaegiLogoImg = loadImage("yaegi.png")
	yaegiLogoOp  = imageOpOrEmpty(yaegiLogoImg)
	gioLogoImg   = loadImage("gio.png")
	gioLogoOp    = imageOpOrEmpty(gioLogoImg)
)

const startingText = `package live

import (
    "gioui.org/layout"
    "gioui.org/widget/material"
    "gioui.org/op"
    "gioui.org/f32"
    "math"
    "image"
)

// store the current rotation offset between frames
var rotation float32

func Layout(gtx layout.Context, theme *material.Theme) layout.Dimensions {
    // Try changing the denomenator of this fraction!
    rotation += math.Pi/50

    // Compute the center of the available area.
    origin := layout.FPt(image.Pt(gtx.Constraints.Max.X/2,gtx.Constraints.Max.Y/2))

    // Spin our drawing around the center.
    op.Affine(f32.Affine2D{}.Rotate(origin, float32(math.Sin(float64(rotation))))).Add(gtx.Ops)

    // Ensure we draw another frame after this one so that animation is smooth.
    op.InvalidateOp{}.Add(gtx.Ops)

    // Draw a word to have something visibly animated.
    return layout.Center.Layout(gtx, func (gtx layout.Context) layout.Dimensions {
        return material.H1(theme, "Hello!").Layout(gtx)
    })
}

`

func squareLogo(gtx C, src image.Image, imgOp *paint.ImageOp) D {
	if src == nil {
		return D{}
	}
	px := gtx.Constraints.Max.Y
	if gtx.Constraints.Max.X < gtx.Constraints.Max.Y {
		px = gtx.Constraints.Max.X
	}
	dps := float32(px) / gtx.Metric.PxPerDp
	scale := dps / float32(imgOp.Size().Y)
	if px != imgOp.Size().X {
		img := image.NewRGBA(image.Rectangle{Max: image.Point{X: px, Y: px}})
		draw.ApproxBiLinear.Scale(img, img.Bounds(), src, src.Bounds(), draw.Src, nil)
		*imgOp = paint.NewImageOp(img)
	}
	return widget.Image{
		Src:   *imgOp,
		Scale: scale,
	}.Layout(gtx)
}

func layoutLogos(gtx C, th *material.Theme) D {
	return layout.Flex{
		Axis:      layout.Vertical,
		Alignment: layout.Middle,
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Center.Layout(gtx, func(gtx C) D {
				return material.Body1(th, "Powered by").Layout(gtx)
			})
		}),
		layout.Flexed(1, func(gtx C) D {
			return layout.Flex{
				Spacing:   layout.SpaceAround,
				Alignment: layout.Middle,
			}.Layout(gtx,
				layout.Flexed(.5, func(gtx C) D {
					return layout.Center.Layout(gtx, func(gtx C) D {
						return squareLogo(gtx, gioLogoImg, &gioLogoOp)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Center.Layout(gtx, func(gtx C) D {
						return material.H3(th, "+").Layout(gtx)
					})
				}),
				layout.Flexed(.5, func(gtx C) D {
					return layout.Center.Layout(gtx, func(gtx C) D {
						return squareLogo(gtx, yaegiLogoImg, &yaegiLogoOp)
					})
				}),
			)
		}),
	)
}

func containsChange(events []widget.EditorEvent) bool {
	for _, e := range events {
		switch e.(type) {
		case widget.ChangeEvent:
			return true
		}
	}
	return false
}

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
			if containsChange(editor.Events()) || first {
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
					inset := layout.UniformInset(unit.Dp(4))
					return inset.Layout(gtx, func(gtx C) D {
						return widget.Border{
							Width: unit.Dp(2),
							Color: th.Fg,
						}.Layout(gtx, func(gtx C) D {
							return inset.Layout(gtx, func(gtx C) D {
								ed := material.Editor(th, &editor, "write layout code here")
								ed.Font.Variant = "Mono"
								return ed.Layout(gtx)
							})
						})
					})

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
						layout.Flexed(.7, func(gtx C) (dims D) {
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
						layout.Flexed(.3, func(gtx C) D {
							return layoutLogos(gtx, th)
						}),
					)
				}),
			)

			e.Frame(gtx.Ops)
			first = false
		}
	}
}
