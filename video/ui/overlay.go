package ui

import (
	"gioui.org/layout"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"golang.org/x/image/colornames"
	"image/color"
)

func (p *Player) drawVideoOverlay(gtx C) {
	paused := p.status == Paused
	switch p.CurrentView() {
	case VideoView:
		if p.uiSeekerWidget.Dragging() || paused {
			if p.uiSeekerWidget.Dragging() {
				cl := color.NRGBA(colornames.Black)
				cl.A = 100
				component.Rect{
					Color: cl,
					Size:  gtx.Constraints.Max,
					Radii: 0,
				}.Layout(gtx)
			}
			layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min.X = 100
				gtx.Constraints.Min.Y = 100
				iconColor := color.NRGBA(colornames.White)
				iconColor.A = 75
				if p.uiSeekerWidget.Dragging() {
					return material.LoaderStyle{
						Color: iconColor,
					}.Layout(gtx)
				}
				return D{}
			})
		}
	}
}
