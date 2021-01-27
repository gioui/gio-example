// SPDX-License-Identifier: Unlicense OR MIT

// GLFW doesn't build on OpenBSD and FreeBSD.
// +build !openbsd,!freebsd,!android,!ios,!js

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
	"math"
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
	"github.com/go-gl/gl/v3.1/gles2"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

// desktopGL is true when the (core, desktop) OpenGL should
// be used, false for OpenGL ES.
const desktopGL = runtime.GOOS == "darwin"

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
	glfw.WindowHint(glfw.ScaleToMonitor, glfw.True)
	glfw.WindowHint(glfw.CocoaRetinaFramebuffer, glfw.True)
	if desktopGL {
		glfw.WindowHint(glfw.ContextVersionMajor, 3)
		glfw.WindowHint(glfw.ContextVersionMinor, 3)
		glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
		glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	} else {
		glfw.WindowHint(glfw.ContextCreationAPI, glfw.EGLContextAPI)
		glfw.WindowHint(glfw.ClientAPI, glfw.OpenGLESAPI)
		glfw.WindowHint(glfw.ContextVersionMajor, 3)
		glfw.WindowHint(glfw.ContextVersionMinor, 0)
	}

	window, err := glfw.CreateWindow(800, 600, "Gio + GLFW", nil, nil)
	if err != nil {
		log.Fatal(err)
	}

	window.MakeContextCurrent()

	if desktopGL {
		err = gl.Init()
	} else {
		err = gles2.Init()
	}
	if err != nil {
		log.Fatalf("gl.Init failed: %v", err)
	}
	if desktopGL {
		// Enable sRGB.
		gl.Enable(gl.FRAMEBUFFER_SRGB)
		// Set up default VBA, required for the forward-compatible core profile.
		var defVBA uint32
		gl.GenVertexArrays(1, &defVBA)
		gl.BindVertexArray(defVBA)
	}

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
		drawOpenGL()
		draw(gtx, th)
		gpu.Collect(sz, gtx.Ops)
		gpu.Frame()
		queue.Frame(gtx.Ops)
		window.SwapBuffers()
	}
}

var (
	button widget.Clickable
	green  float64 = 0.2
)

// drawOpenGL demonstrates the direct use of OpenGL commands
// to draw non-Gio content below the Gio UI.
func drawOpenGL() {
	if desktopGL {
		gl.ClearColor(0, float32(green), 0, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	} else {
		gles2.ClearColor(0, float32(green), 0, 1)
		gles2.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	}
}

// handleCursorEvent handles cursor events not processed by Gio.
func handleCursorEvent(xpos, ypos float64) {
	log.Printf("mouse cursor: (%f,%f)", xpos, ypos)
}

// handleMouseButtonEvent handles mouse button events not processed by Gio.
func handleMouseButtonEvent(button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press {
		green += 0.1
		green, _ = math.Frexp(green)
	}
	log.Printf("mouse button: %v action %v mods %v", button, action, mods)
}

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
		scale := float32(1)
		if runtime.GOOS == "darwin" {
			// macOS cursor positions are not scaled to the underlying framebuffer
			// size when CocoaRetinaFramebuffer is true.
			scale, _ = w.GetContentScale()
		}
		lastPos = f32.Point{X: float32(xpos) * scale, Y: float32(ypos) * scale}
		e := pointer.Event{
			Type:     pointer.Move,
			Position: lastPos,
			Source:   pointer.Mouse,
			Time:     time.Since(beginning),
			Buttons:  btns,
		}
		if !q.Queue(e) {
			handleCursorEvent(xpos, ypos)
		}
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
		e := pointer.Event{
			Type:     typ,
			Source:   pointer.Mouse,
			Time:     time.Since(beginning),
			Position: lastPos,
			Buttons:  btns,
		}
		if !q.Queue(e) {
			handleMouseButtonEvent(button, action, mods)
		}
	})
}
