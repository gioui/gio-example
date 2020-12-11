// SPDX-License-Identifier: Unlicense OR MIT
package main

// A simple Gio program. See https://gioui.org for more information.

import (
	"log"
	"os"
	"strconv"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"git.sr.ht/~whereswaldon/scroll"

	"gioui.org/font/gofont"
)

func main() {
	go func() {
		w := app.NewWindow()
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

var (
	increaseBtn, decreaseBtn widget.Clickable
	length                   int
	list                     layout.List

	indicator scroll.Scrollable
)

type (
	C = layout.Context
	D = layout.Dimensions
)

func loop(w *app.Window) error {
	th := material.NewTheme(gofont.Collection())
	var ops op.Ops
	length = 32
	list.Axis = layout.Vertical
	list.Alignment = layout.Middle
	inset := layout.UniformInset(unit.Dp(4))
	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			if increaseBtn.Clicked() {
				length *= 2
			}
			if decreaseBtn.Clicked() {
				length /= 2
				if length < 1 {
					length = 1
				}
			}
			/*
				Here we check whether the scrollbar experienced user-initiated scrolling
				during the past frame and update the state of the List accordingly
			*/
			if scrolled, progress := indicator.Scrolled(); scrolled {
				list.Position.First = int(float32(length) * progress)
			}
			const third float32 = 1.0 / 3.0
			layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Flex{Alignment: layout.Baseline}.Layout(gtx,
						layout.Flexed(third, func(gtx C) D {
							return inset.Layout(gtx, material.Button(th, &increaseBtn, "Double list length").Layout)
						}),
						layout.Flexed(third, func(gtx C) D {
							return layout.Center.Layout(gtx, func(gtx C) D {
								return material.Body1(th, "Current List Length: "+strconv.Itoa(length)).Layout(gtx)
							})
						}),
						layout.Flexed(third, func(gtx C) D {
							return inset.Layout(gtx, material.Button(th, &decreaseBtn, "Halve list length").Layout)
						}),
					)
				}),
				layout.Flexed(1, func(gtx C) D {
					// Track how many items we are laying out
					var visibleCount int
					dims := list.Layout(gtx, length, func(gtx C, index int) D {
						visibleCount++
						return layout.Center.Layout(gtx, material.H3(th, "List item #"+strconv.Itoa(index)).Layout)
					})
					// Compute (using the heuristic that each item is the same vertical height)
					// the fraction of all items that are currently visible
					visibleFraction := float32(visibleCount) / float32(length)
					// Compute how far through the list items we have scrolled
					scrollDepth := float32(list.Position.First) / float32(length)

					// Lay out the scroll bar
					bar := scroll.DefaultBar(&indicator, scrollDepth, visibleFraction)
					bar.Layout(gtx)
					return dims
				}),
			)
			e.Frame(gtx.Ops)
		}
	}
}
