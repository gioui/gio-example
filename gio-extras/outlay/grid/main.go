package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"strconv"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"git.sr.ht/~whereswaldon/outlay"
)

func main() {
	go func() {
		w := app.NewWindow(
			app.Size(unit.Dp(800), unit.Dp(400)),
			app.Title("Gio layouts"),
		)
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func loop(w *app.Window) error {
	ui := newUI()

	var ops op.Ops
	for e := range w.Events() {
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err

		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			ui.Layout(gtx)
			e.Frame(gtx.Ops)
		}
	}
	return nil
}

type UI struct {
	theme  *material.Theme
	active int
	tabs   []uiTab
	list   layout.List
}

type uiTab struct {
	name  string
	click widget.Clickable
	text  string
	w     func(tab *uiTab, gtx layout.Context) layout.Dimensions
	num   int
	ed    widget.Editor
}

var (
	vWrap = outlay.GridWrap{
		Axis:      layout.Vertical,
		Alignment: layout.End,
	}
	hWrap = outlay.GridWrap{
		Axis:      layout.Horizontal,
		Alignment: layout.End,
	}
	vGrid = outlay.Grid{
		Num:  11,
		Axis: layout.Vertical,
	}
	hGrid = outlay.Grid{
		Num:  11,
		Axis: layout.Horizontal,
	}
)

func newUI() *UI {
	ui := &UI{
		theme: material.NewTheme(gofont.Collection()),
		list: layout.List{
			Axis:      layout.Horizontal,
			Alignment: layout.Baseline,
		},
	}
	ui.tabs = append(ui.tabs,
		uiTab{
			name: "V wrap",
			text: "Lay out items vertically before wrapping to the next column.",
			w: func(tab *uiTab, gtx layout.Context) layout.Dimensions {
				return vWrap.Layout(gtx, tab.num, func(gtx layout.Context, i int) layout.Dimensions {
					s := fmt.Sprintf("item %d", i)
					return material.Body1(ui.theme, s).Layout(gtx)
				})
			},
		},
		uiTab{
			name: "H wrap",
			text: "Lay out items horizontally before wrapping to the next row.",
			w: func(tab *uiTab, gtx layout.Context) layout.Dimensions {
				return hWrap.Layout(gtx, tab.num, func(gtx layout.Context, i int) layout.Dimensions {
					s := fmt.Sprintf("item %d", i)
					return material.Body1(ui.theme, s).Layout(gtx)
				})
			},
		},
		uiTab{
			name: "V grid",
			text: fmt.Sprintf("Lay out %d items vertically before going to the next column.", vGrid.Num),
			w: func(tab *uiTab, gtx layout.Context) layout.Dimensions {
				return vGrid.Layout(gtx, tab.num, func(gtx layout.Context, i int) layout.Dimensions {
					s := fmt.Sprintf("item %d", i)
					return material.Body1(ui.theme, s).Layout(gtx)
				})
			},
		},
		uiTab{
			name: "H grid",
			text: fmt.Sprintf("Lay out %d items horizontally before going to the next row.", hGrid.Num),
			w: func(tab *uiTab, gtx layout.Context) layout.Dimensions {
				return hGrid.Layout(gtx, tab.num, func(gtx layout.Context, i int) layout.Dimensions {
					s := fmt.Sprintf("item %d", i)
					return material.Body1(ui.theme, s).Layout(gtx)
				})
			},
		},
	)
	for i := range ui.tabs {
		tab := &ui.tabs[i]
		tab.ed = widget.Editor{
			SingleLine: true,
			Submit:     true,
		}
		tab.num = 99
		tab.ed.SetText(strconv.Itoa(tab.num))
	}
	return ui
}

func (ui *UI) Layout(gtx layout.Context) layout.Dimensions {
	for i := range ui.tabs {
		for ui.tabs[i].click.Clicked() {
			ui.active = i
		}
	}
	activeTab := &ui.tabs[ui.active]
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return ui.list.Layout(gtx, len(ui.tabs), func(gtx layout.Context, idx int) layout.Dimensions {
				tab := &ui.tabs[idx]
				title := func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(6)).Layout(gtx, material.H6(ui.theme, tab.name).Layout)
				}
				if idx != ui.active {
					return material.Clickable(gtx, &tab.click, title)
				}
				return layout.Stack{}.Layout(gtx,
					layout.Expanded(func(gtx layout.Context) layout.Dimensions {
						clip.UniformRRect(f32.Rectangle{
							Max: layout.FPt(gtx.Constraints.Min),
						}, 0).Add(gtx.Ops)
						paint.Fill(gtx.Ops, color.NRGBA{A: 64})
						return layout.Dimensions{}
					}),
					layout.Stacked(title),
				)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			pt := image.Point{X: gtx.Constraints.Max.X, Y: 4}
			clip.UniformRRect(f32.Rectangle{
				Max: layout.FPt(pt),
			}, 0).Add(gtx.Ops)
			paint.Fill(gtx.Ops, ui.theme.Palette.ContrastBg)
			return layout.Dimensions{Size: pt}
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Stack{}.Layout(gtx,
				layout.Expanded(func(gtx layout.Context) layout.Dimensions {
					clip.UniformRRect(f32.Rectangle{
						Max: layout.FPt(image.Pt(gtx.Constraints.Max.X, gtx.Constraints.Min.Y)),
					}, 0).Add(gtx.Ops)
					paint.Fill(gtx.Ops, color.NRGBA{A: 24})
					return layout.Dimensions{}
				}),
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(4)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						if x, _ := strconv.Atoi(activeTab.ed.Text()); x != activeTab.num {
							activeTab.num = x
						}
						return layout.Flex{
							Alignment: layout.Baseline,
						}.Layout(gtx,
							layout.Rigid(material.Body1(ui.theme, activeTab.text).Layout),
							layout.Rigid(material.Body1(ui.theme, " Num = ").Layout),
							layout.Rigid(material.Editor(ui.theme, &activeTab.ed, "").Layout),
						)
					})
				}),
			)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return activeTab.w(activeTab, gtx)
		}),
	)
}
