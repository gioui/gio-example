// SPDX-License-Identifier: Unlicense OR MIT

package main

// A simple Gio program. See https://gioui.org for more information.

import (
	"fmt"
	"image"
	"image/jpeg"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/io/event"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"gioui.org/font/gofont"
	"gioui.org/x/explorer"
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

type (
	C = layout.Context
	D = layout.Dimensions
)

// ImageResult is the results of trying to open an image. It may
// contain either an error or an image, but not both. The error
// should always be checked first.
type ImageResult struct {
	Error  error
	Format string
	Image  image.Image
}

func loop(w *app.Window) error {
	expl := explorer.NewExplorer(w)
	var openBtn, saveBtn widget.Clickable
	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	imgChan := make(chan ImageResult)
	saveChan := make(chan error)

	events := make(chan event.Event)
	acks := make(chan struct{})

	go func() {
		for {
			ev := w.Event()
			events <- ev
			<-acks
			if _, ok := ev.(app.DestroyEvent); ok {
				return
			}
		}
	}()
	var img ImageResult
	var saveErr error
	var ops op.Ops
	for {
		select {
		case img = <-imgChan:
			w.Invalidate()
		case saveErr = <-saveChan:
			w.Invalidate()
		case e := <-events:
			expl.ListenEvents(e)
			switch e := e.(type) {
			case app.DestroyEvent:
				acks <- struct{}{}
				return e.Err
			case app.FrameEvent:
				gtx := app.NewContext(&ops, e)
				if openBtn.Clicked(gtx) {
					go func() {
						file, err := expl.ChooseFile("png", "jpeg", "jpg")
						if err != nil {
							err = fmt.Errorf("failed opening image file: %w", err)
							imgChan <- ImageResult{Error: err}
							return
						}
						defer file.Close()
						imgData, format, err := image.Decode(file)
						if err != nil {
							err = fmt.Errorf("failed decoding image data: %w", err)
							imgChan <- ImageResult{Error: err}
							return
						}
						imgChan <- ImageResult{Image: imgData, Format: format}
					}()
				}
				if saveBtn.Clicked(gtx) {
					go func(img ImageResult) {
						if img.Error != nil {
							saveChan <- fmt.Errorf("no image loaded, cannot save")
							return
						}
						extension := "jpg"
						switch img.Format {
						case "png":
							extension = "png"
						}
						file, err := expl.CreateFile("file." + extension)
						if err != nil {
							saveChan <- fmt.Errorf("failed exporting image file: %w", err)
							return
						}
						defer func() {
							saveChan <- file.Close()
						}()
						switch extension {
						case "jpg":
							if err := jpeg.Encode(file, img.Image, nil); err != nil {
								saveChan <- fmt.Errorf("failed encoding image file: %w", err)
								return
							}
						case "png":
							if err := png.Encode(file, img.Image); err != nil {
								saveChan <- fmt.Errorf("failed encoding image file: %w", err)
								return
							}
						}
					}(img)
				}
				layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(material.Button(th, &openBtn, "Open Image").Layout),
					layout.Flexed(1, func(gtx C) D {
						if img.Error == nil && img.Image == nil {
							return D{}
						} else if img.Error != nil {
							return material.H6(th, img.Error.Error()).Layout(gtx)
						}

						return widget.Image{
							Src: paint.NewImageOp(img.Image),
							Fit: widget.Contain,
						}.Layout(gtx)
					}),
					layout.Rigid(func(gtx C) D {
						if img.Image == nil {
							gtx = gtx.Disabled()
						}
						return material.Button(th, &saveBtn, "Save Image").Layout(gtx)
					}),
					layout.Rigid(func(gtx C) D {
						if saveErr == nil {
							return D{}
						}
						return material.H6(th, saveErr.Error()).Layout(gtx)
					}),
				)
				e.Frame(gtx.Ops)
			}
			acks <- struct{}{}
		}
	}
}
