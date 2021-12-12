package ui

import (
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"golang.org/x/image/colornames"
	"image/color"
	"time"
)

type View interface {
	Layout(C, ...interface{}) D
}

func (p *Player) IsPathValid() bool {
	return p.isPathValid
}

func (p *Player) FilePath() string {
	return p.filePath
}

func (p *Player) drawStateView(gtx C, args ...interface{}) D {
	var th *material.Theme
	for _, arg := range args {
		switch val := arg.(type) {
		case *material.Theme:
			th = val
		}
	}
	if th == nil {
		th = material.NewTheme(gofont.Collection())
	}

	maxHeight := gtx.Constraints.Max.Y - FooterHeight
	gtx.Constraints.Min.X, gtx.Constraints.Max.X = 600, 600
	gtx.Constraints.Min.Y, gtx.Constraints.Max.Y = maxHeight, maxHeight
	return layout.UniformInset(unit.Dp(16)).Layout(gtx, func(gtx C) D {
		gtx.Constraints.Min = gtx.Constraints.Max
		return p.stateList.Layout(gtx, 10, func(gtx C, index int) D {
			switch index {
			case 0:
				return layout.Inset{Bottom: unit.Dp(32)}.Layout(gtx, func(gtx C) D {
					l := material.H1(th, "Gio Player")
					l.Alignment = text.Middle
					l.Color = th.Bg
					return l.Layout(gtx)
				})
			case 1:
				return layout.Inset{Bottom: unit.Dp(16)}.Layout(gtx, func(gtx C) D {
					fl := layout.Flex{Spacing: layout.SpaceBetween, Alignment: layout.Middle}
					return fl.Layout(gtx,
						layout.Flexed(1, func(gtx C) D {
							return component.Surface(th).Layout(gtx, func(gtx C) D {
								gtx.Constraints.Min.Y = 32
								gtx.Constraints.Min.X = gtx.Constraints.Max.X
								txt := p.FilePath()
								if p.FilePath() == "" {
									txt = "Select video path by clicking button on the right...."
								}
								return layout.UniformInset(unit.Dp(12)).Layout(gtx, func(gtx C) D {
									return material.Label(th, unit.Dp(16), txt).Layout(gtx)
								},
								)
							})
						}),
						layout.Rigid(func(gtx C) D {
							return layout.Spacer{Width: unit.Dp(16)}.Layout(gtx)
						}),
						layout.Rigid(func(gtx C) D {
							icon, _ := widget.NewIcon(icons.FileFolderOpen)
							iconBtn := material.IconButton(th, &p.openFileBtn, icon)
							return iconBtn.Layout(gtx)
						}),
					)
				})
			case 2:
				return layout.Inset{Bottom: unit.Dp(16)}.Layout(gtx, func(gtx C) D {
					fl := layout.Flex{Spacing: layout.SpaceBetween, Alignment: layout.Middle}
					return fl.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							txt := "Audio delay "
							l := material.Label(th, unit.Dp(16), txt)
							l.Color = th.Bg
							return l.Layout(gtx)
						}),
						layout.Rigid(func(gtx C) D {
							return layout.Spacer{Width: unit.Dp(16)}.Layout(gtx)
						}),
						layout.Rigid(func(gtx C) D {
							fl := layout.Flex{}
							return fl.Layout(gtx,
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									icon, _ := widget.NewIcon(icons.ContentRemove)
									return material.IconButton(th, &p.decrementButton, icon).Layout(gtx)
								}),
								layout.Rigid(func(gtx C) D {
									return layout.Spacer{Width: unit.Dp(16)}.Layout(gtx)
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return component.Surface(th).Layout(gtx, func(gtx C) D {
										gtx.Constraints.Min.X = 100
										d := time.Duration(p.audioOffset) * time.Second
										return layout.UniformInset(unit.Dp(12)).Layout(gtx, func(gtx C) D {
											l := material.Label(th, unit.Dp(16), d.String())
											l.Alignment = text.Middle
											return l.Layout(gtx)
										},
										)
									})
								}),
								layout.Rigid(func(gtx C) D {
									return layout.Spacer{Width: unit.Dp(16)}.Layout(gtx)
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									icon, _ := widget.NewIcon(icons.ContentAdd)
									return material.IconButton(th, &p.incrementButton, icon).Layout(gtx)
								}),
							)
						}),
					)
				})

			case 3:
				if !p.IsPathValid() && p.FilePath() != "" {
					return layout.Inset{Bottom: unit.Dp(16)}.Layout(gtx, func(gtx C) D {
						bd := material.Label(th, unit.Dp(16), "Invalid Path")
						bd.Color = color.NRGBA(colornames.Red)
						return bd.Layout(gtx)
					})
				}
				return D{}

			case 4:
				return layout.Inset{Bottom: unit.Dp(16)}.Layout(gtx, func(gtx C) D {
					fl := layout.Flex{Spacing: layout.SpaceBetween, Alignment: layout.Middle}
					return fl.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							txt := "Encoding Status : "
							l := material.Label(th, unit.Dp(16), txt)
							l.Color = th.Bg
							return l.Layout(gtx)
						}),
						layout.Rigid(func(gtx C) D {
							return layout.Spacer{Width: unit.Dp(16)}.Layout(gtx)
						}),
						layout.Rigid(func(gtx C) D {
							s := p.player.EncodingState().Status()
							l := material.Label(th, unit.Dp(16), s.String())
							l.Color = th.Bg
							return l.Layout(gtx)
						}),
					)
				})
			case 5:
				return layout.Inset{Bottom: unit.Dp(16)}.Layout(gtx, func(gtx C) D {
					fl := layout.Flex{Spacing: layout.SpaceBetween, Alignment: layout.Middle}
					return fl.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							txt := "Encoding Message : "
							l := material.Label(th, unit.Dp(16), txt)
							l.Color = th.Bg
							return l.Layout(gtx)
						}),
						layout.Rigid(func(gtx C) D {
							return layout.Spacer{Width: unit.Dp(16)}.Layout(gtx)
						}),
						layout.Rigid(func(gtx C) D {
							s := p.player.EncodingState().Message()
							l := material.Label(th, unit.Dp(16), s)
							l.Color = th.Bg
							return l.Layout(gtx)
						}),
					)
				})
			case 6:
				return layout.Inset{Bottom: unit.Dp(16)}.Layout(gtx, func(gtx C) D {
					fl := layout.Flex{Spacing: layout.SpaceBetween, Alignment: layout.Middle}
					return fl.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							txt := "Decoding Audio Status : "
							l := material.Label(th, unit.Dp(16), txt)
							l.Color = th.Bg
							return l.Layout(gtx)
						}),
						layout.Rigid(func(gtx C) D {
							return layout.Spacer{Width: unit.Dp(16)}.Layout(gtx)
						}),
						layout.Rigid(func(gtx C) D {
							s := p.player.AudioDecodingState().Status()
							l := material.Label(th, unit.Dp(16), s.String())
							l.Color = th.Bg
							return l.Layout(gtx)
						}),
					)
				})
			case 7:
				return layout.Inset{Bottom: unit.Dp(16)}.Layout(gtx, func(gtx C) D {
					fl := layout.Flex{Spacing: layout.SpaceBetween, Alignment: layout.Middle}
					return fl.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							txt := "Decoding Audio Message : "
							l := material.Label(th, unit.Dp(16), txt)
							l.Color = th.Bg
							return l.Layout(gtx)
						}),
						layout.Rigid(func(gtx C) D {
							return layout.Spacer{Width: unit.Dp(16)}.Layout(gtx)
						}),
						layout.Rigid(func(gtx C) D {
							txt := p.player.AudioDecodingState().Message()
							l := material.Label(th, unit.Dp(16), txt)
							l.Color = th.Bg
							return l.Layout(gtx)
						}),
					)
				})
			case 8:
				return layout.Inset{Bottom: unit.Dp(16)}.Layout(gtx, func(gtx C) D {
					fl := layout.Flex{Spacing: layout.SpaceBetween, Alignment: layout.Middle}
					return fl.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							txt := "Decoding Video Status : "
							l := material.Label(th, unit.Dp(16), txt)
							l.Color = th.Bg
							return l.Layout(gtx)
						}),
						layout.Rigid(func(gtx C) D {
							return layout.Spacer{Width: unit.Dp(16)}.Layout(gtx)
						}),
						layout.Rigid(func(gtx C) D {
							s := p.player.VideoDecodingState().Status()
							l := material.Label(th, unit.Dp(16), s.String())
							l.Color = th.Bg
							return l.Layout(gtx)
						}),
					)
				})
			case 9:
				return layout.Inset{Bottom: unit.Dp(16)}.Layout(gtx, func(gtx C) D {
					fl := layout.Flex{Spacing: layout.SpaceBetween, Alignment: layout.Middle}
					return fl.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							txt := "Decoding Video Message : "
							l := material.Label(th, unit.Dp(16), txt)
							l.Color = th.Bg
							return l.Layout(gtx)
						}),
						layout.Rigid(func(gtx C) D {
							return layout.Spacer{Width: unit.Dp(16)}.Layout(gtx)
						}),
						layout.Rigid(func(gtx C) D {
							txt := p.player.VideoDecodingState().Message()
							l := material.Label(th, unit.Dp(16), txt)
							l.Color = th.Bg
							return l.Layout(gtx)
						}),
					)
				})
			}
			return D{}
		})
	})
}
