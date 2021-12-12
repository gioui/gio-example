package ui

import (
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"image"
)

func (p *Player) drawVideoView(gtx C, args ...interface{}) D {
	var img image.Image
	for _, arg := range args {
		switch val := arg.(type) {
		case image.Image:
			img = val
		}
	}
	if img == nil {
		rgb := &image.RGBA{Pix: nil, Stride: 0, Rect: image.Rectangle{}}
		img = rgb.SubImage(rgb.Bounds())
	}
	var imgOps paint.ImageOp
	imgOps = paint.NewImageOp(img)
	return widget.Image{
		Src:      imgOps,
		Fit:      widget.Fill,
		Position: layout.Center,
		Scale:    0,
	}.Layout(gtx)
}
