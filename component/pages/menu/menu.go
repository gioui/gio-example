package menu

import (
	"fmt"
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"

	"gioui.org/example/component/icon"
	page "gioui.org/example/component/pages"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

// Page holds the state for a page demonstrating the features of
// the Menu component.
type Page struct {
	redButton, greenButton, blueButton       widget.Clickable
	balanceButton, accountButton, cartButton widget.Clickable
	leftFillColor                            color.NRGBA
	leftContextArea                          component.ContextArea
	leftMenu, rightMenu                      component.MenuState
	menuInit                                 bool
	menuDemoList                             widget.List
	menuDemoListStates                       []component.ContextArea
	widget.List

	*page.Router
}

// New constructs a Page with the provided router.
func New(router *page.Router) *Page {
	return &Page{
		Router: router,
	}
}

var _ page.Page = &Page{}

func (p *Page) Actions() []component.AppBarAction {
	return []component.AppBarAction{}
}

func (p *Page) Overflow() []component.OverflowAction {
	return []component.OverflowAction{}
}

func (p *Page) NavItem() component.NavItem {
	return component.NavItem{
		Name: "Menu Features",
		Icon: icon.RestaurantMenuIcon,
	}
}

func (p *Page) Layout(gtx C, th *material.Theme) D {
	if !p.menuInit {
		p.leftMenu = component.MenuState{
			Options: []func(gtx C) D{
				func(gtx C) D {
					return layout.Inset{
						Left:  unit.Dp(16),
						Right: unit.Dp(16),
					}.Layout(gtx, material.Body1(th, "Menus support arbitrary widgets.\nThis is just a label!\nHere's a loader:").Layout)
				},
				component.Divider(th).Layout,
				func(gtx C) D {
					return layout.Inset{
						Top:    unit.Dp(4),
						Bottom: unit.Dp(4),
						Left:   unit.Dp(16),
						Right:  unit.Dp(16),
					}.Layout(gtx, func(gtx C) D {
						gtx.Constraints.Max.X = gtx.Dp(unit.Dp(24))
						gtx.Constraints.Max.Y = gtx.Dp(unit.Dp(24))
						return material.Loader(th).Layout(gtx)
					})
				},
				component.SubheadingDivider(th, "Colors").Layout,
				component.MenuItem(th, &p.redButton, "Red").Layout,
				component.MenuItem(th, &p.greenButton, "Green").Layout,
				component.MenuItem(th, &p.blueButton, "Blue").Layout,
			},
		}
		p.rightMenu = component.MenuState{
			Options: []func(gtx C) D{
				func(gtx C) D {
					item := component.MenuItem(th, &p.balanceButton, "Balance")
					item.Icon = icon.AccountBalanceIcon
					item.Hint = component.MenuHintText(th, "Hint")
					return item.Layout(gtx)
				},
				func(gtx C) D {
					item := component.MenuItem(th, &p.accountButton, "Account")
					item.Icon = icon.AccountBoxIcon
					item.Hint = component.MenuHintText(th, "Hint")
					return item.Layout(gtx)
				},
				func(gtx C) D {
					item := component.MenuItem(th, &p.cartButton, "Cart")
					item.Icon = icon.CartIcon
					item.Hint = component.MenuHintText(th, "Hint")
					return item.Layout(gtx)
				},
			},
		}
	}
	if p.redButton.Clicked(gtx) {
		p.leftFillColor = color.NRGBA{R: 200, A: 255}
	}
	if p.greenButton.Clicked(gtx) {
		p.leftFillColor = color.NRGBA{G: 200, A: 255}
	}
	if p.blueButton.Clicked(gtx) {
		p.leftFillColor = color.NRGBA{B: 200, A: 255}
	}
	return layout.Flex{}.Layout(gtx,
		layout.Flexed(.5, func(gtx C) D {
			return widget.Border{
				Color: color.NRGBA{A: 255},
				Width: unit.Dp(2),
			}.Layout(gtx, func(gtx C) D {
				return layout.Stack{}.Layout(gtx,
					layout.Stacked(func(gtx C) D {
						max := image.Pt(gtx.Constraints.Max.X, gtx.Constraints.Max.Y)
						rect := image.Rectangle{
							Max: max,
						}
						paint.FillShape(gtx.Ops, p.leftFillColor, clip.Rect(rect).Op())
						return D{Size: max}
					}),
					layout.Stacked(func(gtx C) D {
						return layout.UniformInset(unit.Dp(12)).Layout(gtx, func(gtx C) D {
							return component.Surface(th).Layout(gtx, func(gtx C) D {
								return layout.UniformInset(unit.Dp(12)).Layout(gtx, material.Body1(th, "Right-click anywhere in this region").Layout)
							})
						})
					}),
					layout.Expanded(func(gtx C) D {
						return p.leftContextArea.Layout(gtx, func(gtx C) D {
							gtx.Constraints.Min = image.Point{}
							return component.Menu(th, &p.leftMenu).Layout(gtx)
						})
					}),
				)
			})
		}),
		layout.Flexed(.5, func(gtx C) D {
			p.menuDemoList.Axis = layout.Vertical
			return material.List(th, &p.menuDemoList).Layout(gtx, 30, func(gtx C, index int) D {
				if len(p.menuDemoListStates) < index+1 {
					p.menuDemoListStates = append(p.menuDemoListStates, component.ContextArea{})
				}
				state := &p.menuDemoListStates[index]
				return layout.Stack{}.Layout(gtx,
					layout.Stacked(func(gtx C) D {
						gtx.Constraints.Min.X = gtx.Constraints.Max.X
						return layout.UniformInset(unit.Dp(8)).Layout(gtx, material.Body1(th, fmt.Sprintf("Item %d", index)).Layout)
					}),
					layout.Expanded(func(gtx C) D {
						return state.Layout(gtx, func(gtx C) D {
							gtx.Constraints.Min.X = 0
							return component.Menu(th, &p.rightMenu).Layout(gtx)
						})
					}),
				)
			})
		}),
	)
}
