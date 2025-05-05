package main

import (
	"log"
	"math"
	"math/rand"
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/example/outlay/fan/playing"
	xwidget "gioui.org/example/outlay/fan/widget"
	"gioui.org/example/outlay/fan/widget/boring"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/outlay"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

func main() {
	go func() {
		w := new(app.Window)
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func genCards(th *material.Theme) []boring.HoverCard {
	cards := []boring.HoverCard{}
	max := 30
	deck := playing.Deck()
	rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})
	for i := range max {
		cards = append(cards, boring.HoverCard{
			CardStyle: boring.CardStyle{
				Card:   deck[i],
				Theme:  th,
				Height: unit.Dp(200),
			},
			HoverState: &xwidget.HoverState{},
		})
	}
	return cards
}

func loop(w *app.Window) error {
	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	fan := outlay.Fan{
		Animation: outlay.Animation{
			Duration: time.Second / 4,
		},
		WidthRadians:  math.Pi,
		OffsetRadians: 2 * math.Pi,
	}
	numCards := widget.Float{}
	numCards.Value = 1.0
	var width, offset, radius widget.Float
	var useRadius widget.Bool
	cardChildren := []outlay.FanItem{}
	cards := genCards(th)
	for i := range cards {
		cardChildren = append(cardChildren, outlay.Item(i == 5, cards[i].Layout))
	}
	var ops op.Ops
	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			for i := range cards {
				cardChildren[i].Elevate = cards[i].Hovering(gtx)
			}
			visibleCards := int(math.Round(float64(numCards.Value*float32(len(cardChildren)-1)))) + 1
			fan.OffsetRadians = offset.Value * 2 * math.Pi
			fan.WidthRadians = width.Value * 2 * math.Pi
			if useRadius.Update(gtx) || radius.Update(gtx) {
				if useRadius.Value {
					r := cards[0].Height * unit.Dp(radius.Value*2)
					fan.HollowRadius = &r
				} else {
					fan.HollowRadius = nil
				}
			}
			layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return material.Body1(th, "1").Layout(gtx)
						}),
						layout.Flexed(1, func(gtx C) D {
							return material.Slider(th, &numCards).Layout(gtx)
						}),
						layout.Rigid(func(gtx C) D {
							return material.Body1(th, "10").Layout(gtx)
						}),
					)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return material.Body1(th, "width 0").Layout(gtx)
						}),
						layout.Flexed(1, func(gtx C) D {
							return material.Slider(th, &width).Layout(gtx)
						}),
						layout.Rigid(func(gtx C) D {
							return material.Body1(th, "2pi").Layout(gtx)
						}),
					)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return material.Body1(th, "offset 0").Layout(gtx)
						}),
						layout.Flexed(1, func(gtx C) D {
							return material.Slider(th, &offset).Layout(gtx)
						}),
						layout.Rigid(func(gtx C) D {
							return material.Body1(th, "2pi").Layout(gtx)
						}),
					)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return material.CheckBox(th, &useRadius, "use").Layout(gtx)
						}),
						layout.Rigid(func(gtx C) D {
							return material.Body1(th, "radius 0%").Layout(gtx)
						}),
						layout.Flexed(1, func(gtx C) D {
							return material.Slider(th, &radius).Layout(gtx)
						}),
						layout.Rigid(func(gtx C) D {
							return material.Body1(th, "200%").Layout(gtx)
						}),
					)
				}),
				layout.Flexed(1, func(gtx C) D {
					return fan.Layout(gtx, cardChildren[:visibleCards]...)
				}),
			)
			e.Frame(gtx.Ops)
		}
	}
}
