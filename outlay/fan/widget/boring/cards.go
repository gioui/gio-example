package boring

import (
	"image/color"
	"math"

	"gioui.org/example/outlay/fan/playing"
	xwidget "gioui.org/example/outlay/fan/widget"
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

type CardPalette struct {
	RedSuit, BlackSuit color.NRGBA
	Border, Background color.NRGBA
}

func (p CardPalette) ColorFor(s playing.Suit) color.NRGBA {
	if s.Color() == playing.Red {
		return p.RedSuit
	}
	return p.BlackSuit
}

var DefaultPalette = &CardPalette{
	RedSuit:    color.NRGBA{R: 0xa0, B: 0x20, A: 0xff},
	BlackSuit:  color.NRGBA{A: 0xff},
	Border:     color.NRGBA{R: 0x80, G: 0x80, B: 0x80, A: 0xff},
	Background: color.NRGBA{R: 0xf0, G: 0xf0, B: 0xf0, A: 0xff},
}

type CardStyle struct {
	*material.Theme
	playing.Card
	Height unit.Dp
	*CardPalette
}

const (
	cardHeightToWidth = 14.0 / 9.0
	cardRadiusToWidth = 1.0 / 16.0
	borderWidth       = 0.005
)

func (c *CardStyle) Palette() *CardPalette {
	if c.CardPalette == nil {
		return DefaultPalette
	}
	return c.CardPalette
}

func (c *CardStyle) Layout(gtx C) D {
	gtx.Constraints.Max.Y = gtx.Dp(c.Height)
	gtx.Constraints.Max.X = int(float32(gtx.Constraints.Max.Y) / cardHeightToWidth)
	outerRadius := float32(gtx.Constraints.Max.X) * cardRadiusToWidth
	innerRadius := (1 - borderWidth) * outerRadius

	borderWidth := c.Height * borderWidth
	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			return Rect{
				Color: c.Palette().Border,
				Size:  layout.FPt(gtx.Constraints.Max),
				Radii: outerRadius,
			}.Layout(gtx)
		}),
		layout.Stacked(func(gtx C) D {
			return layout.UniformInset(borderWidth).Layout(gtx, func(gtx C) D {
				return layout.Stack{}.Layout(gtx,
					layout.Expanded(func(gtx C) D {
						return Rect{
							Color: c.Palette().Background,
							Size:  layout.FPt(gtx.Constraints.Max),
							Radii: innerRadius,
						}.Layout(gtx)
					}),
					layout.Stacked(func(gtx C) D {
						return layout.UniformInset(unit.Dp(2)).Layout(gtx, func(gtx C) D {
							gtx.Constraints.Min = gtx.Constraints.Max
							origin := f32.Point{
								X: float32(gtx.Constraints.Max.X / 2),
								Y: float32(gtx.Constraints.Max.Y / 2),
							}
							layout.Center.Layout(gtx, func(gtx C) D {
								face := material.H1(c.Theme, c.Rank.String())
								face.Color = c.Palette().ColorFor(c.Suit)
								return face.Layout(gtx)
							})
							c.layoutCorner(gtx)
							defer op.Affine(f32.Affine2D{}.Rotate(origin, math.Pi)).Push(gtx.Ops).Pop()
							c.layoutCorner(gtx)

							return D{Size: gtx.Constraints.Max}
						})
					}),
				)
			})
		}),
	)
}

func (c *CardStyle) layoutCorner(gtx layout.Context) layout.Dimensions {
	col := c.Palette().ColorFor(c.Suit)
	return layout.NW.Layout(gtx, func(gtx C) D {
		return layout.UniformInset(unit.Dp(4)).Layout(gtx, func(gtx C) D {
			return layout.Flex{
				Axis:      layout.Vertical,
				Alignment: layout.Middle,
			}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					label := material.H6(c.Theme, c.Rank.String())
					label.Color = col
					return label.Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					label := material.H6(c.Theme, c.Suit.String())
					label.Color = col
					return label.Layout(gtx)
				}),
			)
		})
	})
}

type HoverCard struct {
	CardStyle
	*xwidget.HoverState
}

func (h HoverCard) Layout(gtx C) D {
	dims := h.CardStyle.Layout(gtx)
	gtx.Constraints.Max = dims.Size
	gtx.Constraints.Min = dims.Size
	h.HoverState.Layout(gtx)
	return dims
}
