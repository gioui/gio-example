// SPDX-License-Identifier: Unlicense OR MIT

package main

import (
	"fmt"
	"image"
	"image/color"
	"io"
	"log"
	"os"
	"strings"

	"gioui.org/app"
	"gioui.org/io/event"
	"gioui.org/io/transfer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

const message = "Hello from inside the window."

type demo struct {
	source widget.Draggable
	areas  [2]dropArea
	files  widget.Bool
}

type dropArea struct {
	content string
	mime    string
	drops   int
	active  bool
	err     error
	list    widget.List
}

func main() {
	go func() {
		w := new(app.Window)
		w.Option(app.Title("Drag and drop"), app.Size(800, 480))
		var ops op.Ops
		var d demo
		d.files.Value = true
		theme := material.NewTheme()
		for {
			switch e := w.Event().(type) {
			case app.FrameEvent:
				gtx := app.NewContext(&ops, e)
				d.Layout(gtx, theme)
				e.Frame(&ops)
			case app.DestroyEvent:
				if e.Err != nil {
					log.Fatal(e.Err)
				}
				os.Exit(0)
			}
		}
	}()
	app.Main()
}

func (d *demo) Layout(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	d.source.Type = "text/plain"
	// Offer the text when a drop target requests it.
	if mime, requested := d.source.Update(gtx); requested {
		d.source.Offer(gtx, mime, io.NopCloser(strings.NewReader(message)))
	}
	dragText := func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Min.X = 0
		// Measure the label before painting its background.
		rec := op.Record(gtx.Ops)
		dims := layout.UniformInset(12).Layout(gtx,
			material.Body1(theme, "Drag me: "+message).Layout)
		content := rec.Stop()
		paint.FillShape(gtx.Ops, color.NRGBA{R: 240, G: 242, B: 246, A: 255},
			clip.UniformRRect(image.Rectangle{Max: dims.Size}, gtx.Dp(8)).Op(gtx.Ops))
		content.Add(gtx.Ops)
		return dims
	}
	paint.Fill(gtx.Ops, theme.Bg)
	return layout.UniformInset(16).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return (layout.Flex{Axis: layout.Vertical, Gap: gtx.Dp(16)}).Layout(gtx,
			layout.Rigid(material.H5(theme, "Drag and drop").Layout),
			layout.Rigid(material.Body2(theme,
				"Drag the text below into either area. Enable file drops with the checkbox.").Layout),
			layout.Rigid(material.CheckBox(theme, &d.files, "Accept files").Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return d.source.Layout(gtx, dragText, dragText)
			}),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				axis := layout.Horizontal
				if gtx.Constraints.Max.X < gtx.Dp(600) {
					axis = layout.Vertical
				}
				return (layout.Flex{Axis: axis, Gap: gtx.Dp(16)}).Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return d.areas[0].Layout(gtx, theme, "Area 1", d.files.Value)
					}),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return d.areas[1].Layout(gtx, theme, "Area 2", d.files.Value)
					}),
				)
			}),
		)
	})
}

func (d *dropArea) Layout(gtx layout.Context, theme *material.Theme, title string, files bool) layout.Dimensions {
	filters := []event.Filter{transfer.TargetFilter{Target: d, Type: "text/plain"}}
	if files {
		filters = append(filters, transfer.TargetFilter{Target: d, Type: "text/uri-list"})
	}
	for {
		e, ok := gtx.Event(filters...)
		if !ok {
			break
		}
		switch e := e.(type) {
		case transfer.InitiateEvent:
			d.active = true
		case transfer.CancelEvent:
			d.active = false
		case transfer.DataEvent:
			// Read the data in this frame and close the reader after use.
			reader := e.Open()
			data, err := io.ReadAll(reader)
			closeErr := reader.Close()
			if err == nil {
				err = closeErr
			}
			d.content, d.mime, d.err = string(data), e.Type, err
			d.drops++
			d.list.Position = layout.Position{}
		}
	}
	size := gtx.Constraints.Max
	area := clip.Rect{Max: size}.Push(gtx.Ops)
	defer area.Pop()
	// Keep the whole area droppable, even when empty.
	event.Op(gtx.Ops, d)
	bg := color.NRGBA{R: 235, G: 240, B: 245, A: 255}
	if d.active {
		bg = color.NRGBA{R: 210, G: 235, B: 215, A: 255}
	}
	paint.Fill(gtx.Ops, bg)
	layout.UniformInset(12).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		status := fmt.Sprintf("%d drops received", d.drops)
		if d.mime != "" {
			status += " | " + d.mime
		}
		return (layout.Flex{Axis: layout.Vertical, Gap: gtx.Dp(8)}).Layout(gtx,
			layout.Rigid(material.H6(theme, title).Layout),
			layout.Rigid(material.Body2(theme, status).Layout),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				if d.err != nil {
					label := material.Body2(theme, "Could not read drop: "+d.err.Error())
					label.Color = color.NRGBA{R: 180, G: 30, B: 30, A: 255}
					return label.Layout(gtx)
				}
				if d.drops == 0 {
					text := "Drop text here."
					if files {
						text = "Drop text or files here."
					}
					return material.Body1(theme, text).Layout(gtx)
				}
				lines := strings.Split(strings.TrimRight(d.content, "\r\n"), "\n")
				d.list.Axis = layout.Vertical
				return material.List(theme, &d.list).Layout(gtx, len(lines),
					func(gtx layout.Context, index int) layout.Dimensions {
						return (layout.Inset{Bottom: 12}).Layout(gtx,
							material.Body2(theme, strings.TrimSuffix(lines[index], "\r")).Layout)
					})
			}),
		)
	})
	return layout.Dimensions{Size: size}
}
