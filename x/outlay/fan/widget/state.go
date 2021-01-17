package widget

import (
	"image"

	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

type HoverState struct {
	hovering bool
}

func (c *HoverState) Hovering(gtx C) bool {
	start := c.hovering
	for _, ev := range gtx.Events(c) {
		switch ev := ev.(type) {
		case pointer.Event:
			switch ev.Type {
			case pointer.Enter:
				c.hovering = true
			case pointer.Leave:
				c.hovering = false
			case pointer.Cancel:
				c.hovering = false
			}
		}
	}
	if c.hovering != start {
		op.InvalidateOp{}.Add(gtx.Ops)
	}
	return c.hovering
}

func (c *HoverState) Layout(gtx C) D {
	defer op.Push(gtx.Ops).Pop()
	pointer.Rect(image.Rectangle{Max: gtx.Constraints.Max}).Add(gtx.Ops)
	pointer.InputOp{
		Tag:   c,
		Types: pointer.Enter | pointer.Leave,
	}.Add(gtx.Ops)
	return D{Size: gtx.Constraints.Max}
}
