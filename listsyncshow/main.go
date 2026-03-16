package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

var zcha chan string
var zind int
var zpat string
var zpaa []string

type ImageResult struct {
	Error  error
	Format string
	Image  image.Image
}

func main() {
	go func() {
		fmt.Println("A")
		w1 := new(app.Window)
		w1.Option(app.Title("LISY"))
		w1.Option(app.Size(unit.Dp(600), unit.Dp(600)))
		if err := abc0(w1); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func abc0(w *app.Window) error {
	var xsou widget.Editor
	var xdes widget.Editor
	var xlst widget.Clickable
	var xsyn widget.Clickable
	var xvue widget.Clickable
	var xmsg widget.Editor
	var zsou string
	var zdes string
	zmsg := "LIST prepares csv files of lists of folders and files in your specified source folder."
	zmsg = zmsg + "\nSYNC synchronizes your specified backup folder with your specified source folder."
	zmsg = zmsg + "\nSHOW presents a slideshow of images from your specified source folder and then freeze this window."
	var zops op.Ops
	th := material.NewTheme()
	for {
		switch e := w.Event().(type) {
		case app.FrameEvent:
			gtx := app.NewContext(&zops, e)
			if xlst.Clicked(gtx) {
				zlin := ""
				zlrd := ""
				zlrf := ""
				ztag := time.Now().Format("06-01-02_15-04")
				ztop := strings.TrimSpace(xsou.Text())
				err := filepath.Walk(ztop,
					func(path string, info os.FileInfo, err error) error {
						if err != nil {
							return err
						}
						zlin = strings.Replace(path, ztop, "", -1) + "," + strconv.Itoa(int(info.Size())) + "," + info.ModTime().Format("06-01-02_15-04")
						fmt.Println(zlin)
						if info.IsDir() {
							zlrd = zlrd + zlin + "\r\n"
						} else {
							zlrf = zlrf + zlin + "\r\n"
						}
						return nil
					})
				if err != nil {
					log.Println(err)
				} else {
					os.WriteFile(ztop+"\\ListDir_"+ztag+".csv", []byte("RelPath,Bytes,ModYmdhm\r\n"+zlrd), 0666)
					os.WriteFile(ztop+"\\ListFil_"+ztag+".csv", []byte("RelPath,Bytes,ModYmdhm\r\n"+zlrf), 0666)
				}
				zmsg = "See: " + ztop + "\\ListDir(Fil)_" + ztag + ".csv"
				fmt.Println(zmsg)
				zmsg = zmsg + "\n" + zlrf
			}
			if xvue.Clicked(gtx) {
				abc1(strings.TrimSpace(xsou.Text()))
			}
			if xsyn.Clicked(gtx) {
				zsou = strings.TrimSpace(xsou.Text())
				zdes = strings.TrimSpace(xdes.Text())
				if len(xdes.Text()) > 0 {
					zcmd := exec.Command("ROBOCOPY", zsou, zdes, "/MIR")
					zout, err := zcmd.CombinedOutput()
					if err != nil {
						zmsg = "Error: " + string(err.Error())
					} else {
						zmsg = string(zout)
					}
					fmt.Println(zmsg)
				}
				xmsg.SetText(zmsg)
			}
			layout.Flex{
				Axis:    layout.Vertical,
				Spacing: layout.SpaceStart,
			}.Layout(
				gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return material.Label(th, unit.Sp(10), "Specify your source folder, for List or Sync or View.").Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					ed := material.Editor(th, &xsou, "Source folder (for List Sync View)")
					xsou.SingleLine = false
					xsou.Alignment = text.Start
					margins := layout.Inset{
						Top:    unit.Dp(1),
						Bottom: unit.Dp(1),
						Left:   unit.Dp(10),
						Right:  unit.Dp(10),
					}
					border := widget.Border{
						Color:        color.NRGBA{R: 204, G: 204, B: 204, A: 255},
						CornerRadius: unit.Dp(5),
						Width:        unit.Dp(1),
					}
					return margins.Layout(gtx,
						func(gtx layout.Context) layout.Dimensions {
							return border.Layout(gtx, ed.Layout)
						},
					)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return material.Label(th, unit.Sp(10), "Specify your backup folder, for Sync.").Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					ed := material.Editor(th, &xdes, "Backup folder (for Sync)")
					xdes.SingleLine = false
					xdes.Alignment = text.Start
					margins := layout.Inset{
						Top:    unit.Dp(1),
						Bottom: unit.Dp(1),
						Left:   unit.Dp(10),
						Right:  unit.Dp(10),
					}
					border := widget.Border{
						Color:        color.NRGBA{R: 204, G: 204, B: 204, A: 255},
						CornerRadius: unit.Dp(5),
						Width:        unit.Dp(1),
					}
					return margins.Layout(gtx,
						func(gtx layout.Context) layout.Dimensions {
							return border.Layout(gtx, ed.Layout)
						},
					)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return material.Label(th, unit.Sp(10), "Click List or Sync or View button.").Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{}.Layout(gtx,
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							return material.Button(th, &xlst, "LIST").Layout(gtx)
						}),
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							return material.Button(th, &xsyn, "SYNC").Layout(gtx)
						}),
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							return material.Button(th, &xvue, "SHOW").Layout(gtx)
						}),
					)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return material.Label(th, unit.Sp(10), "Note the following.").Layout(gtx)
				}),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					ed := material.Editor(th, &xmsg, zmsg)
					xmsg.SingleLine = false
					xmsg.Alignment = text.Start
					margins := layout.Inset{
						Top:    unit.Dp(1),
						Bottom: unit.Dp(1),
						Left:   unit.Dp(10),
						Right:  unit.Dp(10),
					}
					border := widget.Border{
						Color:        color.NRGBA{R: 204, G: 204, B: 204, A: 255},
						CornerRadius: unit.Dp(5),
						Width:        unit.Dp(1),
					}
					return margins.Layout(gtx,
						func(gtx layout.Context) layout.Dimensions {
							return border.Layout(gtx, ed.Layout)
						},
					)
				}),
			)
			e.Frame(gtx.Ops)
		case app.DestroyEvent:
			return e.Err
		}
	}
}

func abc1(ztop string) {
	zpat = ""
	filepath.Walk(ztop,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			} else {
				if info.IsDir() {
					fmt.Println("Dir: " + path)
				} else {
					if strings.HasSuffix(path, ".jpg") || strings.HasSuffix(path, ".png") {
						fmt.Println("File: " + path)
						zpat = zpat + "," + path
					}
				}
			}
			return nil
		})
	zpaa = strings.Split(zpat, ",")
	fmt.Println(len(zpaa))
	go func() {
		w2 := new(app.Window)
		w2.Option(app.Title("LisyImages"))
		if err := abc2(w2); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	zcha = make(chan string, 1)
	zind = 0
	go func() {
		for {
			if zind < len(zpaa)-2 {
				time.Sleep(time.Second * 2)
				zcha <- "LisyImages"
			}
		}
	}()
	app.Main()
}

func abc2(w *app.Window) error {
	var ops op.Ops
	th := material.NewTheme()
	var openBtn widget.Clickable
	for {
		switch e := w.Event().(type) {
		case app.FrameEvent:
			select {
			case ztit := <-zcha:
				if zind > len(zpaa)-1 {
					// w.Close()
					return nil
				} else {
					zind = zind + 1
					fmt.Println(zind)
					gtx := app.NewContext(&ops, e)
					w.Option(app.Title(ztit))
					file, err := os.OpenFile(strings.Split(zpat, ",")[zind], 0, 0)
					if err != nil {
						err = fmt.Errorf("failed opening image file: %w", err)
					}
					defer file.Close()
					imgData, format, err := image.Decode(file)
					if err != nil {
						err = fmt.Errorf("failed decoding image data: %w", err)
					}
					img := ImageResult{Image: imgData, Format: format}
					layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return material.Button(th, &openBtn, zpaa[zind]).Layout(gtx)
						}),
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							return widget.Image{
								Src: paint.NewImageOp(img.Image),
								Fit: widget.Contain,
							}.Layout(gtx)
						}),
					)
					e.Frame(gtx.Ops)
				}
			}
		}
	}
}
