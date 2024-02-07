package widget

import (
	"image"

	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
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
	for {
		ev, ok := gtx.Event(pointer.Filter{
			Target: c,
			Kinds:  pointer.Enter | pointer.Leave,
		})
		if !ok {
			break
		}
		switch ev := ev.(type) {
		case pointer.Event:
			switch ev.Kind {
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
		gtx.Execute(op.InvalidateCmd{})
	}
	return c.hovering
}

func (c *HoverState) Layout(gtx C) D {
	defer clip.Rect(image.Rectangle{Max: gtx.Constraints.Max}).Push(gtx.Ops).Pop()
	event.Op(gtx.Ops, c)
	return D{Size: gtx.Constraints.Max}
}
