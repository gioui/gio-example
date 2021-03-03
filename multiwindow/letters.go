// SPDX-License-Identifier: Unlicense OR MIT

package main

import (
	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// Letters displays a clickable list of text items that open a new window.
type Letters struct {
	win *Window
	log *Log

	items []*LetterListItem
	list  layout.List
}

// NewLetters creates a new letters view with the provided log.
func NewLetters(log *Log) *Letters {
	view := &Letters{
		log:  log,
		list: layout.List{Axis: layout.Vertical},
	}
	for text := 'a'; text <= 'z'; text++ {
		view.items = append(view.items, &LetterListItem{Text: string(text)})
	}
	return view
}

// Run implements Window.Run method.
func (v *Letters) Run(w *Window) error {
	v.win = w
	return WidgetView(v.Layout).Run(w)
}

// Layout handles drawing the letters view.
func (v *Letters) Layout(gtx layout.Context) layout.Dimensions {
	th := v.win.App.Theme
	return v.list.Layout(gtx, len(v.items), func(gtx layout.Context, index int) layout.Dimensions {
		item := v.items[index]
		for item.Click.Clicked() {
			v.log.Printf("opening %s view", item.Text)

			bigText := material.H1(th, item.Text)
			size := bigText.TextSize
			size.V *= 2
			v.win.App.NewWindow(item.Text,
				WidgetView(func(gtx layout.Context) layout.Dimensions {
					return layout.Center.Layout(gtx, bigText.Layout)
				}),
				app.Size(size, size),
			)
		}
		return material.Button(th, &item.Click, item.Text).Layout(gtx)
	})
}

type LetterListItem struct {
	Text  string
	Click widget.Clickable
}
