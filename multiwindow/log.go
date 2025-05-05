// SPDX-License-Identifier: Unlicense OR MIT

package main

import (
	"fmt"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/event"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// Log shows a list of strings.
type Log struct {
	addLine chan string
	lines   []string

	close widget.Clickable
	list  widget.List
}

// NewLog crates a new log view.
func NewLog() *Log {
	return &Log{
		addLine: make(chan string, 100),
		list:    widget.List{List: layout.List{Axis: layout.Vertical}},
	}
}

// Printf adds a new line to the log.
func (log *Log) Printf(format string, args ...any) {
	s := fmt.Sprintf(format, args...)

	// ensure that this logging does not block.
	select {
	case log.addLine <- s:
	default:
	}
}

// Run handles window loop for the log.
func (log *Log) Run(w *Window) error {
	var ops op.Ops

	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))

	go func() {
		<-w.App.Context.Done()
		w.Perform(system.ActionClose)
	}()

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
	for {
		select {
		// listen to new lines from Printf and add them to our lines.
		case line := <-log.addLine:
			log.lines = append(log.lines, line)
			w.Invalidate()
		case e := <-events:
			switch e := e.(type) {
			case app.DestroyEvent:
				acks <- struct{}{}
				return e.Err
			case app.FrameEvent:
				gtx := app.NewContext(&ops, e)
				log.Layout(w, th, gtx)
				e.Frame(gtx.Ops)
			}
			acks <- struct{}{}
		}
	}
}

// Layout displays the log with a close button.
func (log *Log) Layout(w *Window, th *material.Theme, gtx layout.Context) {
	// This is here to demonstrate programmatic closing of a window,
	// however it's probably better to use OS close button instead.
	for log.close.Clicked(gtx) {
		w.Window.Perform(system.ActionClose)
	}

	layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(material.Button(th, &log.close, "Close").Layout),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return material.List(th, &log.list).Layout(gtx, len(log.lines), func(gtx layout.Context, i int) layout.Dimensions {
				return material.Body1(th, log.lines[i]).Layout(gtx)
			})
		}),
	)
}
