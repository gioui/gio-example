package about

import (
	"gioui.org/io/clipboard"
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"

	alo "gioui.org/example/component/applayout"
	"gioui.org/example/component/icon"
	page "gioui.org/example/component/pages"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

// Page holds the state for a page demonstrating the features of
// the AppBar component.
type Page struct {
	eliasCopyButton, chrisCopyButtonGH, chrisCopyButtonLP widget.Clickable
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
		Name: "About this library",
		Icon: icon.OtherIcon,
	}
}

const (
	sponsorEliasURL          = "https://github.com/sponsors/eliasnaur"
	sponsorChrisURLGitHub    = "https://github.com/sponsors/whereswaldon"
	sponsorChrisURLLiberapay = "https://liberapay.com/whereswaldon/"
)

func (p *Page) Layout(gtx C, th *material.Theme) D {
	p.List.Axis = layout.Vertical
	return material.List(th, &p.List).Layout(gtx, 1, func(gtx C, _ int) D {
		return layout.Flex{
			Alignment: layout.Middle,
			Axis:      layout.Vertical,
		}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return alo.DefaultInset.Layout(gtx, material.Body1(th, `This library implements material design components from https://material.io using https://gioui.org.

If you like this library and work like it, please consider sponsoring Elias and/or Chris!`).Layout)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return alo.DetailRow{}.Layout(gtx,
					material.Body1(th, "Elias Naur can be sponsored on GitHub at "+sponsorEliasURL).Layout,
					func(gtx C) D {
						if p.eliasCopyButton.Clicked() {
							clipboard.WriteOp{
								Text: sponsorEliasURL,
							}.Add(gtx.Ops)
						}
						return material.Button(th, &p.eliasCopyButton, "Copy Sponsorship URL").Layout(gtx)
					})
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return alo.DetailRow{}.Layout(gtx,
					material.Body1(th, "Chris Waldon can be sponsored on GitHub at "+sponsorChrisURLGitHub+" and on Liberapay at "+sponsorChrisURLLiberapay).Layout,

					func(gtx C) D {
						if p.chrisCopyButtonGH.Clicked() {
							clipboard.WriteOp{Text: sponsorChrisURLGitHub}.Add(gtx.Ops)
						}
						if p.chrisCopyButtonLP.Clicked() {
							clipboard.WriteOp{Text: sponsorChrisURLLiberapay}.Add(gtx.Ops)
						}
						return alo.DefaultInset.Layout(gtx, func(gtx C) D {
							return layout.Flex{}.Layout(gtx,
								layout.Flexed(.5, material.Button(th, &p.chrisCopyButtonGH, "Copy GitHub URL").Layout),
								layout.Flexed(.5, material.Button(th, &p.chrisCopyButtonLP, "Copy Liberapay URL").Layout),
							)
						})
					})
			}),
		)
	})
}
