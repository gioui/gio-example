// SPDX-License-Identifier: Unlicense OR MIT

package main

import (
	"image"
	"math/rand"
)

// Board implements game of life logic.
type Board struct {
	// Size is the count of cells in a particular dimension.
	Size image.Point
	// Cells contains the alive or dead cells.
	Cells []byte

	// buffer is used to avoid reallocating a new cells
	// slice for every update.
	buffer []byte
}

// NewBoard returns a new game of life with the defined size.
func NewBoard(size image.Point) *Board {
	return &Board{
		Size:   size,
		Cells:  make([]byte, size.X*size.Y),
		buffer: make([]byte, size.X*size.Y),
	}
}

// Randomize randomizes each cell state.
func (b *Board) Randomize() {
	rand.Read(b.Cells)
	for i, v := range b.Cells {
		if v < 0x30 {
			b.Cells[i] = 1
		} else {
			b.Cells[i] = 0
		}
	}
}

// Pt returns the coordinate given a index in b.Cells.
func (b *Board) Pt(i int) image.Point {
	x, y := i%b.Size.X, i/b.Size.Y
	return image.Point{X: x, Y: y}
}

// At returns the b.Cells index, given a wrapped coordinate.
func (b *Board) At(c image.Point) int {
	if c.X < 0 {
		c.X += b.Size.X
	}
	if c.X >= b.Size.X {
		c.X -= b.Size.X
	}
	if c.Y < 0 {
		c.Y += b.Size.Y
	}
	if c.Y >= b.Size.Y {
		c.Y -= b.Size.Y
	}
	return b.Size.Y*c.Y + c.X
}

// SetWithoutWrap sets a cell to alive.
func (b *Board) SetWithoutWrap(c image.Point) {
	if !c.In(image.Rectangle{Max: b.Size}) {
		return
	}

	b.Cells[b.At(c)] = 1
}

// Advance advances the board state by 1.
func (b *Board) Advance() {
	next, cur := b.buffer, b.Cells
	defer func() { b.Cells, b.buffer = next, cur }()

	for i := range next {
		next[i] = 0
	}

	for y := range b.Size.Y {
		for x := range b.Size.X {
			var t byte
			t += cur[b.At(image.Pt(x-1, y-1))]
			t += cur[b.At(image.Pt(x+0, y-1))]
			t += cur[b.At(image.Pt(x+1, y-1))]
			t += cur[b.At(image.Pt(x-1, y+0))]
			t += cur[b.At(image.Pt(x+1, y+0))]
			t += cur[b.At(image.Pt(x-1, y+1))]
			t += cur[b.At(image.Pt(x+0, y+1))]
			t += cur[b.At(image.Pt(x+1, y+1))]

			// Any live cell with fewer than two live neighbours dies, as if by underpopulation.
			// Any live cell with two or three live neighbours lives on to the next generation.
			// Any live cell with more than three live neighbours dies, as if by overpopulation.
			// Any dead cell with exactly three live neighbours becomes a live cell, as if by reproduction.

			p := b.At(image.Pt(x, y))
			switch {
			case t < 2:
				t = 0
			case t == 2:
				t = cur[p]
			case t == 3:
				t = 1
			case t > 3:
				t = 0
			}

			next[p] = t
		}
	}
}
