// SPDX-License-Identifier: Unlicense OR MIT

// GLFW doesn't build on OpenBSD and FreeBSD.
// +build !openbsd,!freebsd,!windows,!android,!ios,!js

// The glfw example demonstrates integration of Gio into a foreign
// windowing and rendering library, in this case GLFW
// (https://www.glfw.org).
//
// See the go-glfw package for installation of the native
// dependencies:
//
// https://github.com/go-gl/glfw
package main

import (
	"image"
	"log"
	"runtime"
	"time"

	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/gpu"
	giogl "gioui.org/gpu/gl"
	"gioui.org/io/pointer"
	"gioui.org/io/router"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func main() {
	// Required by the OpenGL threading model.
	runtime.LockOSThread()

	err := glfw.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer glfw.Terminate()
	// Gio assumes a sRGB backbuffer.
	glfw.WindowHint(glfw.SRGBCapable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.CocoaRetinaFramebuffer, glfw.True)

	window, err := glfw.CreateWindow(800, 600, "Gio + GLFW", nil, nil)
	if err != nil {
		log.Fatal(err)
	}

	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		log.Fatal(err)
	}
	// Enable sRGB.
	gl.Enable(gl.FRAMEBUFFER_SRGB)
	// Set up default VBA, required for the forward-compatible core profile.
	var defVBA uint32
	gl.GenVertexArrays(1, &defVBA)
	gl.BindVertexArray(defVBA)

	var queue router.Router
	var ops op.Ops
	th := material.NewTheme(gofont.Collection())
	backend, err := giogl.NewBackend(nil)
	if err != nil {
		log.Fatal(err)
	}
	gpu, err := gpu.New(backend)
	if err != nil {
		log.Fatal(err)
	}

	registerCallbacks(window, &queue)
	for !window.ShouldClose() {
		glfw.PollEvents()
		scale, _ := window.GetContentScale()
		width, height := window.GetFramebufferSize()
		sz := image.Point{X: width, Y: height}
		ops.Reset()
		gtx := layout.Context{
			Ops:   &ops,
			Now:   time.Now(),
			Queue: &queue,
			Metric: unit.Metric{
				PxPerDp: scale,
				PxPerSp: scale,
			},
			Constraints: layout.Exact(sz),
		}
		draw(gtx, th)
		gpu.Collect(sz, gtx.Ops)
		gpu.Frame()
		queue.Frame(gtx.Ops)
		window.SwapBuffers()
	}
}

var button widget.Clickable

func draw(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return layout.Center.Layout(gtx,
		material.Button(th, &button, "Button").Layout,
	)
}

func registerCallbacks(window *glfw.Window, q *router.Router) {
	var btns pointer.Buttons
	beginning := time.Now()
	var lastPos f32.Point
	window.SetCursorPosCallback(func(w *glfw.Window, xpos float64, ypos float64) {
		scale, _ := w.GetContentScale()
		lastPos = f32.Point{X: float32(xpos) * scale, Y: float32(ypos) * scale}
		q.Add(pointer.Event{
			Type:     pointer.Move,
			Position: lastPos,
			Source:   pointer.Mouse,
			Time:     time.Since(beginning),
			Buttons:  btns,
		})
	})
	window.SetMouseButtonCallback(func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
		var btn pointer.Buttons
		switch button {
		case glfw.MouseButton1:
			btn = pointer.ButtonLeft
		case glfw.MouseButton2:
			btn = pointer.ButtonRight
		case glfw.MouseButton3:
			btn = pointer.ButtonMiddle
		}
		var typ pointer.Type
		switch action {
		case glfw.Release:
			typ = pointer.Release
			btns &^= btn
		case glfw.Press:
			typ = pointer.Press
			btns |= btn
		}
		q.Add(pointer.Event{
			Type:     typ,
			Source:   pointer.Mouse,
			Time:     time.Since(beginning),
			Position: lastPos,
			Buttons:  btns,
		})
	})
}
