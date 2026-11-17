package main

import (
	"fmt"
	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/io/pointer"
	"gioui.org/io/profile"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"gioui.org/x/gamepad"
	"image"
	"image/color"
	"math"
	"os"
)

var (
	theme = material.NewTheme(gofont.Collection())
    disableCache = false
	disableGamepad = false
)

func main() {
	w := app.NewWindow(app.Title("Gamepad Tester"))

	gamepads := gamepad.NewGamepad(w)
	ops := new(op.Ops)

	page := Page{
		Controllers: [4]Controller{
			{Controller: gamepads.Controllers[0]},
			{Controller: gamepads.Controllers[1]},
			{Controller: gamepads.Controllers[2]},
			{Controller: gamepads.Controllers[3]},
		},
	}

	go func() {
		for evt := range w.Events() {
			if !disableGamepad {
				gamepads.ListenEvents(evt)
			}

			switch evt := evt.(type) {
			case pointer.Event:
				if evt.Buttons.Contain(pointer.ButtonSecondary) && evt.Type == pointer.Press {
					disableGamepad = !disableGamepad
				}
			case system.FrameEvent:
				gtx := layout.NewContext(ops, evt)

				if gamepads.Controllers[0].Buttons.Start.Pressed {
					disableCache = true
				}
				if gamepads.Controllers[0].Buttons.Back.Pressed {
					disableCache = false
				}

				page.Layout(gtx)

				op.InvalidateOp{}.Add(gtx.Ops)
				evt.Frame(ops)
			case system.DestroyEvent:
				if evt.Err != nil {
					panic(evt.Err)
				}
				os.Exit(0)
			}
		}
	}()

	app.Main()
}

type Page struct {
	Controllers [4]Controller
	profile     profile.Event
}

func (p *Page) Layout(gtx layout.Context) {
	p.Controllers[0].Layout(gtx)

	for _, e := range gtx.Events(p) {
		if e, ok := e.(profile.Event); ok {
			p.profile = e
		}
	}
	profile.Op{Tag: p}.Add(gtx.Ops)

	layout.NE.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Top: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			txt := fmt.Sprintf("cache: %v (back/start) | gamepad: %v (right-click) | m: %s", !disableCache, !disableGamepad, p.profile.Timings)
			lbl := material.H6(theme, txt)
			lbl.Font.Variant = "Mono"
			return lbl.Layout(gtx)
		})
	})
}

type Controller struct {
	Color      color.NRGBA
	Controller *gamepad.Controller
}

func (c Controller) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Max.X = gtx.Px(unit.Dp(300))
		gtx.Constraints.Max.Y = gtx.Px(unit.Dp(160))
		gtx.Constraints.Min = gtx.Constraints.Max

		return layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceAround}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceAround}.Layout(gtx,
					layout.Rigid(Trigger{
						Color:  color.NRGBA{A: 255},
						Button: c.Controller.Buttons.LT,
					}.Layout),
					layout.Rigid(layout.Spacer{Width: unit.Dp(70)}.Layout),
					layout.Rigid(Trigger{
						Color:  color.NRGBA{A: 255},
						Button: c.Controller.Buttons.RT,
					}.Layout),
				)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceAround}.Layout(gtx,
					layout.Rigid(Shoulder{
						Color:  color.NRGBA{A: 255},
						Button: c.Controller.Buttons.LB,
					}.Layout),
					layout.Rigid(layout.Spacer{Width: unit.Dp(70)}.Layout),
					layout.Rigid(Shoulder{
						Color:  color.NRGBA{A: 255},
						Button: c.Controller.Buttons.RB,
					}.Layout),
				)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceAround}.Layout(gtx,
					layout.Rigid(Joystick{
						Color:    color.NRGBA{A: 255},
						Position: c.Controller.Joysticks.LeftThumb,
						Button:   c.Controller.Buttons.LeftThumb,
					}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Max.X, gtx.Constraints.Max.Y = gtx.Px(unit.Dp(70)), gtx.Px(unit.Dp(50))
						gtx.Constraints.Min = gtx.Constraints.Max

						return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceSides, Alignment: layout.Middle}.Layout(gtx,
							layout.Rigid(OptionsButtons{
								Color:  c.Color,
								Button: c.Controller.Buttons.Back,
							}.Layout),
							layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
							layout.Rigid(OptionsButtons{
								Color:  c.Color,
								Button: c.Controller.Buttons.Start,
							}.Layout),
						)
					}),

					layout.Rigid(Buttons{
						Color: color.NRGBA{A: 255},
						Up:    c.Controller.Buttons.Y,
						Down:  c.Controller.Buttons.A,
						Left:  c.Controller.Buttons.X,
						Right: c.Controller.Buttons.B,
					}.Layout),
				)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceSides}.Layout(gtx,
					layout.Rigid(Buttons{
						Color: color.NRGBA{A: 255},
						Up:    c.Controller.Buttons.Up,
						Down:  c.Controller.Buttons.Down,
						Left:  c.Controller.Buttons.Left,
						Right: c.Controller.Buttons.Right,
					}.Layout),
					layout.Rigid(layout.Spacer{Width: unit.Dp(50)}.Layout),
					layout.Rigid(Joystick{
						Color:    color.NRGBA{A: 255},
						Position: c.Controller.Joysticks.RightThumb,
						Button:   c.Controller.Buttons.RightThumb,
					}.Layout),
				)
			}),
		)
	})
}

type OptionsButtons struct {
	Color  color.NRGBA
	Button gamepad.Button
}

func (o OptionsButtons) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Max.X, gtx.Constraints.Max.Y = gtx.Px(unit.Dp(30)), gtx.Px(unit.Dp(20))

		paint.FillShape(gtx.Ops, opacity(o.Color, uint8(128)), stroke(circle(gtx, gtx), gtx, unit.Dp(3)))
		paint.FillShape(gtx.Ops, opacity(o.Color, o.Button.Force), infill(circle(gtx, gtx)))

		return layout.Dimensions{Size: gtx.Constraints.Max}
	})
}

type Shoulder struct {
	Color  color.NRGBA
	Button gamepad.Button
}

func (o Shoulder) Layout(gtx layout.Context) layout.Dimensions {
	gtx.Constraints.Max.X, gtx.Constraints.Max.Y = gtx.Px(unit.Dp(50)), gtx.Px(unit.Dp(12))

	paint.FillShape(gtx.Ops, opacity(o.Color, uint8(128)), stroke(circle(gtx, gtx), gtx, unit.Dp(3)))
	paint.FillShape(gtx.Ops, opacity(o.Color, o.Button.Force), infill(circle(gtx, gtx)))

	return layout.Dimensions{Size: gtx.Constraints.Max}
}

type Trigger struct {
	Color  color.NRGBA
	Button gamepad.Button
}

func (o Trigger) Layout(gtx layout.Context) layout.Dimensions {
	gtx.Constraints.Max.X, gtx.Constraints.Max.Y = gtx.Px(unit.Dp(50)), gtx.Px(unit.Dp(14))

	paint.FillShape(gtx.Ops, opacity(o.Color, uint8(128)), stroke(circle(gtx, gtx), gtx, unit.Dp(3)))
	paint.FillShape(gtx.Ops, opacity(o.Color, o.Button.Force), infill(circle(gtx, gtx)))

	return layout.Dimensions{Size: gtx.Constraints.Max}
}

type Button struct {
	Color  color.NRGBA
	Button gamepad.Button
}

func (b Button) Layout(gtx layout.Context) layout.Dimensions {
	return layout.UniformInset(unit.Dp(2)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		paint.FillShape(gtx.Ops, opacity(b.Color, uint8(128)), stroke(circle(gtx, gtx), gtx, unit.Dp(3)))
		paint.FillShape(gtx.Ops, opacity(b.Color, b.Button.Force), infill(circle(gtx, gtx)))
		return layout.Dimensions{Size: gtx.Constraints.Max}
	})
}

type Buttons struct {
	Color                 color.NRGBA
	Up, Down, Left, Right gamepad.Button
}

func (b Buttons) Layout(gtx layout.Context) layout.Dimensions {
	gtx.Constraints.Max.X, gtx.Constraints.Max.Y = gtx.Px(unit.Dp(50)), gtx.Px(unit.Dp(50))
	gtx.Constraints.Min = gtx.Constraints.Max

	paint.FillShape(gtx.Ops, opacity(b.Color, uint8(128)), stroke(circle(gtx, gtx), gtx, unit.Dp(3)))

	defer op.Affine(f32.Affine2D{}.Rotate(getSizeFloat(gtx).Div(2), math.Phi/2)).Push(gtx.Ops).Pop()

	return layout.UniformInset(unit.Dp(5)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceSides}.Layout(gtx,
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceSides}.Layout(gtx,
					layout.Flexed(1, Button{Color: b.Color, Button: b.Up}.Layout),
					layout.Flexed(1, Button{Color: b.Color, Button: b.Right}.Layout),
				)
			}),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceSides}.Layout(gtx,
					layout.Flexed(1, Button{Color: b.Color, Button: b.Left}.Layout),
					layout.Flexed(1, Button{Color: b.Color, Button: b.Down}.Layout),
				)
			}),
		)
	})
}

type Joystick struct {
	Color    color.NRGBA
	Button   gamepad.Button
	Position f32.Point
}

func (j Joystick) Layout(gtx layout.Context) layout.Dimensions {
	gtx.Constraints.Max.X, gtx.Constraints.Max.Y = gtx.Px(unit.Dp(50)), gtx.Px(unit.Dp(50))

	paint.FillShape(gtx.Ops, opacity(j.Color, uint8(128)), stroke(circle(gtx, gtx), gtx, unit.Dp(3)))
	paint.FillShape(gtx.Ops, opacity(j.Color, j.Button.Force), stroke(circle(gtx, gtx), gtx, unit.Dp(3)))
	return layout.UniformInset(unit.Dp(5)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		defer op.Offset(j.Position.Mul(8)).Push(gtx.Ops).Pop()

		paint.FillShape(gtx.Ops, j.Color, infill(circle(gtx, gtx)))
		return layout.Dimensions{Size: gtx.Constraints.Max}
	})
}

/* HELPER */

type Opacity interface{
	~float32 | ~float64 | ~uint8 | ~int
}

func opacity[O Opacity](c color.NRGBA, opacity O) color.NRGBA {
	cc := *(&c)
	switch opacity := (interface{})(opacity).(type) {
	case float32:
		cc.A = uint8(opacity * 255)
	case float64:
		cc.A = uint8(opacity * 255)
	case uint8:
		cc.A = opacity
	}
	return cc
}

type Size interface {
	image.Point | f32.Point | layout.Constraints | f32.Rectangle | layout.Context
}

func getSizeFloat[S Size](size S) f32.Point {
	switch size := (interface{})(size).(type) {
	case image.Point:
		return f32.Pt(float32(size.X), float32(size.Y))
	case layout.Constraints:
		return f32.Pt(float32(size.Max.X), float32(size.Max.Y))
	case f32.Point:
		return f32.Pt(size.X, size.Y)
	case f32.Rectangle:
		return f32.Pt(size.Max.X, size.Min.Y)
	case layout.Context:
		return f32.Pt(float32(size.Constraints.Max.X), float32(size.Constraints.Max.Y))
	}
	return f32.Pt(0, 0)
}

type Unit interface{
~int | ~int64 | unit.Value | ~float64 | ~float32
}

func getUnitFloat[U Unit](metric unit.Metric, u U) float32 {
	switch radius := (interface{})(u).(type) {
	case unit.Value:
		return float32(metric.Px(radius))
	case int:
		return float32(radius)
	case int64:
		return float32(radius)
	case float64:
		return float32(radius)
	case float32:
		return radius
	default:
		return 0
	}
}

var rectCache = map[clip.RRect]clip.PathSpec{}

func rect[U Unit, S Size](gtx layout.Context, size S, radius U) clip.PathSpec {
	c := clip.UniformRRect(f32.Rectangle{Max: getSizeFloat(size)}, getUnitFloat(gtx.Metric, radius))
	if disableCache {
		return c.Path(gtx.Ops)
	}
	cached, ok := rectCache[c]
	if !ok {
		cached = c.Path(new(op.Ops))
		rectCache[c] = cached
	}
	return cached
}

func circle[S Size](gtx layout.Context, size S) clip.PathSpec {
	ss := getSizeFloat(size)
	return rect(gtx, size, ss.Y/2)
}

var strokeCache = map[clip.Stroke]clip.Op{}

func stroke[U Unit](path clip.PathSpec, gtx layout.Context, width U) clip.Op {
	c := clip.Stroke{Path: path, Width: getUnitFloat(gtx.Metric, width)}
	if disableCache {
		return c.Op()
	}
	cached, ok := strokeCache[c]
	if !ok {
		fmt.Println("not cached")
		cached = c.Op()
		strokeCache[c] = cached
	}
	return cached
}

var infillCache = map[clip.Outline]clip.Op{}

func infill(path clip.PathSpec) clip.Op {
	c := clip.Outline{Path: path}
	if disableCache {
		return c.Op()
	}
	cached, ok := infillCache[c]
	if !ok {
		cached = c.Op()
		infillCache[c] = cached
	}
	return cached
}
