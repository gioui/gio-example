package discloser

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"

	"gioui.org/example/component/icon"
	page "gioui.org/example/component/pages"
)

// TreeNode is a simple tree implementation that holds both
// display data and the state for Discloser widgets. In
// practice, you'll often want to separate the state from
// the data being presented.
type TreeNode struct {
	Text     string
	Children []TreeNode
	component.DiscloserState
}

type (
	C = layout.Context
	D = layout.Dimensions
)

// Page holds the state for a page demonstrating the features of
// the AppBar component.
type Page struct {
	TreeNode
	widget.List
	*page.Router
	CustomDiscloserState component.DiscloserState
}

// New constructs a Page with the provided router.
func New(router *page.Router) *Page {
	return &Page{
		Router: router,
		TreeNode: TreeNode{
			Text: "Expand Me",
			Children: []TreeNode{
				{
					Text: "Disclosers can be (expand me)...",
					Children: []TreeNode{
						{
							Text: "...nested to arbitrary depths.",
						},
						{
							Text: "There are also types available to customize the look and feel of the discloser:",
							Children: []TreeNode{
								{
									Text: "• DiscloserStyle lets you provide your own control instead of the default triangle used here.",
								},
								{
									Text: "• DiscloserArrowStyle lets you alter the presentation of the triangle used here, like changing its color, size, left/right anchoring, or margin.",
								},
							},
						},
					},
				},
			},
		},
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
		Name: "Disclosers",
		Icon: icon.VisibilityIcon,
	}
}

// LayoutTreeNode recursively lays out a tree of widgets described by
// TreeNodes.
func (p *Page) LayoutTreeNode(gtx C, th *material.Theme, tn *TreeNode) D {
	if len(tn.Children) == 0 {
		return layout.UniformInset(unit.Dp(2)).Layout(gtx,
			material.Body1(th, tn.Text).Layout)
	}
	children := make([]layout.FlexChild, 0, len(tn.Children))
	for i := range tn.Children {
		child := &tn.Children[i]
		children = append(children, layout.Rigid(
			func(gtx C) D {
				return p.LayoutTreeNode(gtx, th, child)
			}))
	}
	return component.SimpleDiscloser(th, &tn.DiscloserState).Layout(gtx,
		material.Body1(th, tn.Text).Layout,
		func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
		})
}

// LayoutCustomDiscloser demonstrates how to create a custom control for
// a discloser.
func (p *Page) LayoutCustomDiscloser(gtx C, th *material.Theme) D {
	return component.Discloser(th, &p.CustomDiscloserState).Layout(gtx,
		func(gtx C) D {
			var l material.LabelStyle
			l = material.Body1(th, "+")
			if p.CustomDiscloserState.Visible() {
				l.Text = "-"
			}
			l.Font.Typeface = "Go Mono"
			return layout.UniformInset(unit.Dp(2)).Layout(gtx, l.Layout)
		},
		material.Body1(th, "Custom Control").Layout,
		material.Body2(th, "This control only took 9 lines of code.").Layout,
	)
}

func (p *Page) Layout(gtx C, th *material.Theme) D {
	p.List.Axis = layout.Vertical
	return material.List(th, &p.List).Layout(gtx, 2, func(gtx C, index int) D {
		return layout.UniformInset(unit.Dp(4)).Layout(gtx, func(gtx C) D {
			if index == 0 {
				return p.LayoutTreeNode(gtx, th, &p.TreeNode)
			}
			return p.LayoutCustomDiscloser(gtx, th)
		})
	})
}
