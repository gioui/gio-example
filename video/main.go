package main

import (
	"gioui.org/app"
	"gioui.org/example/video/ui"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"log"
	"os"
	"time"
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
	player := ui.Player{Window: window}

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
				player.Layout(gtx)
				e.Frame(gtx.Ops)
			case key.Event:
				if e.Name == "→" && e.Modifiers == 0x0 && e.State == 0x0 {
					player.RightKeyPressed = true
					player.RightKeyLastUpdated = time.Now().UnixMilli()
				}
				if e.Name == "→" && e.Modifiers == 0x0 && e.State == 0x1 {
					player.RightKeyPressed = false
					player.RightKeyLastUpdated = time.Now().UnixMilli()
				}
				if e.Name == "↑" && e.Modifiers == 0x0 && e.State == 0x0 {
					player.UpKeyPressed = true
					player.UpKeyLastUpdated = time.Now().UnixMilli()
				}
				if e.Name == "↑" && e.Modifiers == 0x0 && e.State == 0x1 {
					player.UpKeyPressed = false
					player.RightKeyLastUpdated = time.Now().UnixMilli()
				}
				if e.Name == "←" && e.Modifiers == 0x0 && e.State == 0x0 {
					player.LeftKeyPressed = true
					player.LeftKeyLastUpdated = time.Now().UnixMilli()
				}
				if e.Name == "←" && e.Modifiers == 0x0 && e.State == 0x1 {
					player.LeftKeyPressed = false
					player.RightKeyLastUpdated = time.Now().UnixMilli()
				}
				if e.Name == "↓" && e.Modifiers == 0x0 && e.State == 0x0 {
					player.DownKeyPressed = true
					player.DownKeyLastUpdated = time.Now().UnixMilli()
				}
				if e.Name == "↓" && e.Modifiers == 0x0 && e.State == 0x1 {
					player.DownKeyPressed = false
					player.RightKeyLastUpdated = time.Now().UnixMilli()
				}
				if e.Name == "Space" && e.Modifiers == 0x0 && e.State == 0x0 {
					player.SpaceKeyPressed = true
					player.SpaceKeyLastUpdated = time.Now().UnixMilli()
				}
				if e.Name == "Space" && e.Modifiers == 0x0 && e.State == 0x1 {
					player.SpaceKeyPressed = false
					player.SpaceKeyLastUpdated = time.Now().UnixMilli()
				}
				if e.Name == "⎋" && e.Modifiers == 0x0 && e.State == 0x0 {
					player.EscKeyPressed = true
					player.EscKeyLastUpdated = time.Now().UnixMilli()
				}
				if e.Name == "⎋" && e.Modifiers == 0x0 && e.State == 0x1 {
					player.EscKeyPressed = false
					player.EscKeyLastUpdated = time.Now().UnixMilli()
				}
			}
		}
	}
}
