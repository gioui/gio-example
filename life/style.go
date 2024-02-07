// SPDX-License-Identifier: Unlicense OR MIT

package main

import (
	"image"
	"image/color"

	"gioui.org/f32" // f32 is used for shape calculations.
	"gioui.org/io/event"
	"gioui.org/io/pointer" // system is used for system events (e.g. closing the window).
	"gioui.org/layout"     // layout is used for layouting widgets.

	// op is used for recording different operations.
	"gioui.org/op/clip"  // clip is used to draw the cell shape.
	"gioui.org/op/paint" // paint is used to paint the cells.
)

// BoardStyle draws Board with rectangles.
type BoardStyle struct {
	CellSizePx int
	*Board
}

// Layout draws the Board and accepts input for adding alive cells.
func (board BoardStyle) Layout(gtx layout.Context) layout.Dimensions {
	// Calculate the board size based on the cell size in pixels.
	size := board.Size.Mul(board.CellSizePx)
	gtx.Constraints = layout.Exact(size)

	// Handle any input from a pointer.
	for {
		ev, ok := gtx.Event(pointer.Filter{
			Target: board.Board,
			Kinds:  pointer.Drag,
		})
		if !ok {
			break
		}
		if ev, ok := ev.(pointer.Event); ok {
			p := image.Pt(int(ev.Position.X), int(ev.Position.Y))
			// Calculate the board coordinate given a cursor position.
			p = p.Div(board.CellSizePx)
			board.SetWithoutWrap(p)
		}
	}
	// Register to listen for pointer Drag events.
	pr := clip.Rect(image.Rectangle{Max: size}).Push(gtx.Ops)
	event.Op(gtx.Ops, board.Board)
	pr.Pop()

	cellSize := float32(board.CellSizePx)

	// Draw a shape for each alive cell.
	var p clip.Path
	p.Begin(gtx.Ops)
	for i, v := range board.Cells {
		if v == 0 {
			continue
		}

		c := layout.FPt(board.Pt(i).Mul(board.CellSizePx))
		p.MoveTo(f32.Pt(c.X, c.Y))
		p.LineTo(f32.Pt(c.X+cellSize, c.Y))
		p.LineTo(f32.Pt(c.X+cellSize, c.Y+cellSize))
		p.LineTo(f32.Pt(c.X, c.Y+cellSize))
		p.Close()
	}
	defer clip.Outline{Path: p.End()}.Op().Push(gtx.Ops).Pop()

	// Paint the shape with a black color.
	paint.ColorOp{Color: color.NRGBA{A: 0xFF}}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	return layout.Dimensions{Size: size}
}
