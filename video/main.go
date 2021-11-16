package main

import (
	"gioui.org/app"
	video "gioui.org/example/video/player"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"log"
	"os"
	"path/filepath"
)

func main() {
	go func() {
		w := app.NewWindow(app.Title("Gio Player"))
		if err := loop(w); err != nil {
			log.Fatalln(err)
		}
		os.Exit(0)
	}()
	app.Main()
}
func loop(window *app.Window) (err error) {
	videoPath, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}
	videoPath = filepath.Join(videoPath, "gio.mp4")
	p := video.Player{Filepath: videoPath}
	go p.Start()

	var ops op.Ops

	for {
		select {
		case e := <-window.Events():
			switch e := e.(type) {
			case system.DestroyEvent:
				return e.Err
			case *system.CommandEvent:
				switch e.Type {
				case system.CommandBack:
					log.Fatalln(e)
				}
			case system.FrameEvent:
				gtx := layout.NewContext(&ops, e)
				p.Layout(gtx)
				e.Frame(gtx.Ops)
			}
		}
	}
}