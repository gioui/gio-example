package main

import (
	"log"
	"math"
	"math/rand"
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/example/x/outlay/fan/playing"
	xwidget "gioui.org/example/x/outlay/fan/widget"
	"gioui.org/example/x/outlay/fan/widget/boring"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
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
		w := app.NewWindow()
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
	for i := 0; i < max; i++ {
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
	th := material.NewTheme(gofont.Collection())
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
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			for i := range cards {
				cardChildren[i].Elevate = cards[i].Hovering(gtx)
			}
			visibleCards := int(math.Round(float64(numCards.Value*float32(len(cardChildren)-1)))) + 1
			fan.OffsetRadians = offset.Value * 2 * math.Pi
			fan.WidthRadians = width.Value * 2 * math.Pi
			if useRadius.Changed() || radius.Changed() {
				if useRadius.Value {
					r := cards[0].Height.Scale(radius.Value * 2)
					fan.HollowRadius = &r
				} else {
					fan.HollowRadius = nil
				}
			}
			layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Flex{}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return material.Body1(th, "1").Layout(gtx)
						}),
						layout.Flexed(1, func(gtx C) D {
							return material.Slider(th, &numCards, 0.0, 1.0).Layout(gtx)
						}),
						layout.Rigid(func(gtx C) D {
							return material.Body1(th, "10").Layout(gtx)
						}),
					)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Flex{}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return material.Body1(th, "width 0").Layout(gtx)
						}),
						layout.Flexed(1, func(gtx C) D {
							return material.Slider(th, &width, 0.0, 1.0).Layout(gtx)
						}),
						layout.Rigid(func(gtx C) D {
							return material.Body1(th, "2pi").Layout(gtx)
						}),
					)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Flex{}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return material.Body1(th, "offset 0").Layout(gtx)
						}),
						layout.Flexed(1, func(gtx C) D {
							return material.Slider(th, &offset, 0.0, 1.0).Layout(gtx)
						}),
						layout.Rigid(func(gtx C) D {
							return material.Body1(th, "2pi").Layout(gtx)
						}),
					)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Flex{}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return material.CheckBox(th, &useRadius, "use").Layout(gtx)
						}),
						layout.Rigid(func(gtx C) D {
							return material.Body1(th, "radius 0%").Layout(gtx)
						}),
						layout.Flexed(1, func(gtx C) D {
							return material.Slider(th, &radius, 0.0, 1.0).Layout(gtx)
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
