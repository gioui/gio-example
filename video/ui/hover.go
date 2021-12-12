package ui

import (
	"gioui.org/io/pointer"
	"image"
)

func (p *Player) drawHover(gtx C) D {
	gtx.Constraints.Min = gtx.Constraints.Max
	defer pointer.Rect(image.Rectangle{Max: gtx.Constraints.Max}).Push(gtx.Ops).Pop()
	defer p.screenClickable.Layout(gtx)
	pointer.InputOp{
		Tag:   &p.lastHoveredTime,
		Types: pointer.Move,
	}.Add(gtx.Ops)
	return D{Size: gtx.Constraints.Max}
}
