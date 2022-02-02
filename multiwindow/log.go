// SPDX-License-Identifier: Unlicense OR MIT

package main

import (
	"fmt"

	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
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
func (log *Log) Printf(format string, args ...interface{}) {
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

	th := material.NewTheme(gofont.Collection())

	applicationClose := w.App.Context.Done()
	for {
		select {
		case <-applicationClose:
			return nil
		// listen to new lines from Printf and add them to our lines.
		case line := <-log.addLine:
			log.lines = append(log.lines, line)
			w.Invalidate()
		case e := <-w.Events():
			switch e := e.(type) {
			case system.DestroyEvent:
				return e.Err
			case system.FrameEvent:
				gtx := layout.NewContext(&ops, e)
				log.Layout(w, th, gtx)
				e.Frame(gtx.Ops)
			}
		}
	}
}

// Layout displays the log with a close button.
func (log *Log) Layout(w *Window, th *material.Theme, gtx layout.Context) {
	// This is here to demonstrate programmatic closing of a window,
	// however it's probably better to use OS close button instead.
	for log.close.Clicked() {
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
