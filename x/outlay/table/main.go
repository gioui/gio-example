package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"strconv"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/outlay"
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
	var table = uiTable{
		theme: material.NewTheme(gofont.Collection()),
		Table: outlay.Table{
			CellSize: func(m unit.Metric, x, y int) image.Point {
				return image.Pt(m.Px(unit.Dp(50)), m.Px(unit.Dp(30)))
			},
		},
	}

	var ops op.Ops
	for e := range w.Events() {
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err

		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			table.Layout(gtx)
			e.Frame(gtx.Ops)
		}
	}
	return nil
}

type uiTable struct {
	theme    *material.Theme
	xed, yed widget.Editor
	cells    []cell
	outlay.Table
}

type cell struct {
	click   widget.Clickable
	clicked bool
}

func (t *uiTable) Layout(gtx layout.Context) layout.Dimensions {
	th := t.theme
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis:    layout.Horizontal,
				Spacing: layout.SpaceSides,
			}.Layout(gtx,
				layout.Rigid(material.Body1(th, "Number of columns: ").Layout),
				layout.Rigid(material.Editor(th, &t.xed, "columns").Layout),
			)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis:    layout.Horizontal,
				Spacing: layout.SpaceSides,
			}.Layout(gtx,
				layout.Rigid(material.Body1(th, "Number of rows: ").Layout),
				layout.Rigid(material.Editor(th, &t.yed, "rows").Layout),
			)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis:    layout.Horizontal,
				Spacing: layout.SpaceSides,
			}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					var selected int
					for i := range t.cells {
						if t.cells[i].clicked {
							selected++
						}
					}
					var txt string
					switch selected {
					case 0:
						txt = "Click cells to select them."
					case 1:
						txt = fmt.Sprint("1 cell selected. Click again to unselect.")
					default:
						txt = fmt.Sprintf("%d cells selected. Click again to unselect.", selected)
					}
					return material.Body1(th, txt).Layout(gtx)
				}),
			)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			xn, _ := strconv.Atoi(t.xed.Text())
			yn, _ := strconv.Atoi(t.yed.Text())
			t.cells = growCells(t.cells, xn*yn)
			return t.Table.Layout(gtx, xn, yn, func(gtx layout.Context, x, y int) layout.Dimensions {
				c := &t.cells[x+y*xn]
				return c.Layout(gtx, th, x, y)
			})
		}),
	)
}

func (c *cell) Layout(gtx layout.Context, th *material.Theme, x, y int) layout.Dimensions {
	defer op.Save(gtx.Ops).Load()
	var txt string
	if y < 0 {
		txt = fmt.Sprintf("item %d", x)
	} else {
		txt = fmt.Sprintf("%dx%d", x, y)
	}
	macro := op.Record(gtx.Ops)
	dims := material.Clickable(gtx, &c.click, func(gtx layout.Context) layout.Dimensions {
		return layout.Center.Layout(gtx, material.Body1(th, txt).Layout)
	})
	call := macro.Stop()

	if c.click.Clicked() {
		c.clicked = !c.clicked
	}
	clip.Rect{Max: dims.Size}.Add(gtx.Ops)
	col := color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	if c.clicked {
		col = color.NRGBA{R: 128, G: 128, B: 128, A: 255}
	}
	paint.Fill(gtx.Ops, col)
	call.Add(gtx.Ops)
	return dims
}

func growCells(cells []cell, n int) []cell {
	if len(cells) < n {
		cells = append(cells, make([]cell, n-len(cells))...)
	}
	return cells[:n]
}
