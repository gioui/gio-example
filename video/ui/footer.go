package ui

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"time"
)

func getFormattedTime(val time.Duration) string {
	d := val
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

const FooterHeight = 120

func (p *Player) drawFooter(gtx C) D {
	gtx.Constraints.Min.Y, gtx.Constraints.Max.Y = FooterHeight, FooterHeight
	ins := layout.Inset{Top: unit.Dp(16), Right: unit.Dp(8), Bottom: unit.Dp(16), Left: unit.Dp(8)}
	fl := layout.Flex{Axis: layout.Vertical, Alignment: layout.End}
	rgd1 := layout.Rigid(func(gtx C) D {
		return layout.Inset{
			Right: unit.Dp(16),
			Left:  unit.Dp(16),
		}.Layout(gtx,
			func(gtx C) D {
				return layout.Flex{
					Alignment: layout.Middle,
				}.Layout(gtx,
					layout.Flexed(1, p.drawSlider),
					layout.Rigid(func(gtx C) D {
						return layout.Spacer{Width: unit.Dp(16)}.Layout(gtx)
					}),
					layout.Rigid(func(gtx C) D {
						value := time.Duration(0)
						//if p.videoDuration.Seconds() != 0 {
						if p.uiSeekerWidget.Value != 0 {
							totalSeconds := p.uiSeekerWidget.Value
							value = time.Duration(totalSeconds) * time.Second
						}
						text := getFormattedTime(value)
						l := material.Label(p.th, unit.Dp(16), text)
						if p.CurrentView() == StateView {
							l.Color = p.th.Bg
						}
						return l.Layout(gtx)
					}),
					layout.Rigid(func(gtx C) D {
						d := time.Duration(p.player.EndTime()) * time.Second
						text := " / " + getFormattedTime(d)
						l := material.Label(p.th, unit.Dp(16), text)
						if p.CurrentView() == StateView {
							l.Color = p.th.Bg
						}
						return l.Layout(gtx)
					}),
				)
			},
		)
	})
	rgd2 := layout.Rigid(func(gtx C) D {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		return layout.Flex{
			Spacing:   layout.SpaceSides,
			Alignment: layout.Middle,
			WeightSum: 0,
		}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Center.Layout(gtx, func(gtx C) D {
					icon, _ := widget.NewIcon(icons.AVStop)
					return material.IconButton(p.th, &p.stopBtn, icon).Layout(gtx)
				})
			}),
			layout.Rigid(func(gtx C) D {
				return layout.Spacer{Width: unit.Dp(32)}.Layout(gtx)
			}),
			layout.Rigid(func(gtx C) D {
				return layout.Center.Layout(gtx, func(gtx C) D {
					icon, _ := widget.NewIcon(icons.AVFastRewind)
					return material.IconButton(p.th, &p.rewindButton, icon).Layout(gtx)
				})
			}),
			layout.Rigid(func(gtx C) D {
				return layout.Spacer{Width: unit.Dp(32)}.Layout(gtx)
			}),
			layout.Rigid(func(gtx C) D {
				return layout.Center.Layout(gtx, func(gtx C) D {
					icon, _ := widget.NewIcon(icons.AVPlayArrow)
					btn := &p.playBtn
					if p.status == Playing {
						icon, _ = widget.NewIcon(icons.AVPause)
						btn = &p.pauseBtn
					}
					return material.IconButton(p.th, btn, icon).Layout(gtx)
				})
			}),
			layout.Rigid(func(gtx C) D {
				return layout.Spacer{Width: unit.Dp(32)}.Layout(gtx)
			}),
			layout.Rigid(func(gtx C) D {
				return layout.Center.Layout(gtx, func(gtx C) D {
					icon, _ := widget.NewIcon(icons.AVFastForward)
					return material.IconButton(p.th, &p.forwardButton, icon).Layout(gtx)
				})
			}),
			layout.Rigid(func(gtx C) D {
				return layout.Spacer{Width: unit.Dp(32)}.Layout(gtx)
			}),
			layout.Rigid(func(gtx C) D {
				return layout.Center.Layout(gtx, func(gtx C) D {
					icon, _ := widget.NewIcon(icons.ActionHelp)
					if p.CurrentView() == StateView {
						icon, _ = widget.NewIcon(icons.AVMovie)
					}
					return material.IconButton(p.th, &p.showStateBtn, icon).Layout(gtx)
				})
			}),
		)
	})

	wgt := func(gtx C) D {
		return fl.Layout(gtx, rgd1, rgd2)
	}
	return ins.Layout(gtx, wgt)
}
