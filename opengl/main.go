// SPDX-License-Identifier: Unlicense OR MIT

//go:build darwin || windows || linux
// +build darwin windows linux

// This program demonstrates the use of a custom OpenGL ES context with
// app.Window. It is similar to the GLFW example, but uses Gio's window
// implementation instead of the one in GLFW.
//
// The example runs on Linux using the normal EGL and X11 libraries, so
// no additional libraries need to be installed. However, it must be
// build with -tags nowayland until app.ViewEvent is implemented for
// Wayland.
//
// The example runs on macOS and Windows using ANGLE:
//
// $ CGO_CFLAGS=-I<path-to-ANGLE>/include CGO_LDFLAGS=-L<path-to-angle-libraries> go build -o opengl.exe ./opengl
//
// You'll need the ANGLE libraries (EGL and GLESv2) in your library search path. On macOS:
//
// $ DYLD_LIBRARY_PATH=<path-to-ANGLE-libraries> ./opengl.exe
package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"runtime"
	"strings"
	"unsafe"

	"gioui.org/app"
	"gioui.org/gpu"
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"gioui.org/font/gofont"
)

/*
#cgo CFLAGS: -DEGL_NO_X11
#cgo LDFLAGS: -lEGL -lGLESv2

#include <EGL/egl.h>
#include <GLES3/gl3.h>
#define EGL_EGLEXT_PROTOTYPES
#include <EGL/eglext.h>

*/
import "C"

type eglContext struct {
	view    C.EGLNativeWindowType
	disp    C.EGLDisplay
	ctx     C.EGLContext
	surf    C.EGLSurface
	cleanup func()
}

func main() {
	go func() {
		// Set CustomRenderer so we can provide our own rendering context.
		w := new(app.Window)
		w.Option(app.CustomRenderer(true))
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

var btnScreenshot widget.Clickable

func discardContext(w *app.Window, gioGPU gpu.GPU, eglCtx *eglContext) error {
	var err error
	w.Run(func() {
		if gioGPU != nil {
			gioGPU.Release()
		}
		if eglCtx != nil {
			if ok := C.eglMakeCurrent(eglCtx.disp, eglCtx.surf, eglCtx.surf, eglCtx.ctx); ok != C.EGL_TRUE {
				err = fmt.Errorf("eglMakeCurrent failed (%#x)", C.eglGetError())
			} else {
				eglCtx.Release()
			}
		}
	})
	return err
}

func recreateContext(w *app.Window, ve app.ViewEvent, gioGPU gpu.GPU, eglCtx *eglContext, size image.Point) (gpu.GPU, *eglContext, error) {
	err := discardContext(w, gioGPU, eglCtx)
	gioGPU = nil
	eglCtx = nil
	if err != nil {
		return nil, nil, err
	}
	w.Run(func() {
		eglCtx, err = createContext(ve, size)
		if err != nil {
			err = fmt.Errorf("failed creating context: %w", err)
		}
	})
	if err != nil {
		return nil, nil, err
	}
	if ok := C.eglMakeCurrent(eglCtx.disp, eglCtx.surf, eglCtx.surf, eglCtx.ctx); ok != C.EGL_TRUE {
		return gioGPU, eglCtx, fmt.Errorf("eglMakeCurrent failed (%#x)", C.eglGetError())
	}
	glGetString := func(e C.GLenum) string {
		return C.GoString((*C.char)(unsafe.Pointer(C.glGetString(e))))
	}
	fmt.Printf("GL_VERSION: %s\nGL_RENDERER: %s\n", glGetString(C.GL_VERSION), glGetString(C.GL_RENDERER))
	gioGPU, err = gpu.New(gpu.OpenGL{ES: true, Shared: true})
	return gioGPU, eglCtx, err
}

func loop(w *app.Window) (windowErr error) {
	// eglMakeCurrent binds a context to an operating system thread. Prevent Go from switching thread.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	var ops op.Ops
	var (
		ctx    *eglContext
		gioCtx gpu.GPU
		ve     app.ViewEvent
		size   image.Point
	)

	for {
		switch e := w.Event().(type) {
		case app.ViewEvent:
			ve = e
			err := discardContext(w, gioCtx, ctx)
			windowErr = errors.Join(windowErr, err)
			gioCtx = nil
			ctx = nil
		case app.DestroyEvent:
			windowErr = errors.Join(windowErr, e.Err)
			return windowErr
		case app.FrameEvent:
			resized := e.Size != size
			size = e.Size
			if ve.Valid() && size != (image.Point{}) && gioCtx == nil {
				var err error
				gioCtx, ctx, err = recreateContext(w, ve, gioCtx, ctx, size)
				windowErr = errors.Join(windowErr, err)
			} else if ve.Valid() && resized {
				ctx.resizeSurface(ve, size)
			}
			if windowErr != nil {
				w.Perform(system.ActionClose)
				continue
			}
			// Build ops.
			gtx := app.NewContext(&ops, e)
			// Catch pointer events not hitting UI.
			types := pointer.Move | pointer.Press | pointer.Release
			event.Op(gtx.Ops, w)
			for {
				e, ok := gtx.Event(pointer.Filter{
					Target: w,
					Kinds:  types,
				})
				if !ok {
					break
				}
				log.Println("Event:", e)
			}
			if btnScreenshot.Clicked(gtx) {
				if err := screenshot(gioCtx, e.Size, gtx.Ops); err != nil {
					log.Printf("failed taking screenshot: %v", err)
				}
			}

			drawUI(th, gtx)
			// Trigger window resize detection in ANGLE.
			C.eglWaitClient()
			// Draw custom OpenGL content.
			drawGL()

			// Render drawing ops.
			if err := gioCtx.Frame(gtx.Ops, gpu.OpenGLRenderTarget{}, e.Size); err != nil {
				windowErr = errors.Join(windowErr, fmt.Errorf("render failed: %v", err))
				continue
			}

			if ok := C.eglSwapBuffers(ctx.disp, ctx.surf); ok != C.EGL_TRUE {
				windowErr = errors.Join(windowErr, fmt.Errorf("swap failed: 0x%x", C.eglGetError()))
				continue
			}

			// Process non-drawing ops.
			e.Frame(gtx.Ops)
		}
	}
}

func screenshot(ctx gpu.GPU, size image.Point, ops *op.Ops) error {
	var tex C.GLuint
	C.glGenTextures(1, &tex)
	defer C.glDeleteTextures(1, &tex)
	C.glBindTexture(C.GL_TEXTURE_2D, tex)
	C.glTexImage2D(C.GL_TEXTURE_2D, 0, C.GL_RGBA, C.GLint(size.X), C.GLint(size.Y), 0, C.GL_RGBA, C.GL_UNSIGNED_BYTE, nil)
	var fbo C.GLuint
	C.glGenFramebuffers(1, &fbo)
	defer C.glDeleteFramebuffers(1, &fbo)
	C.glBindFramebuffer(C.GL_FRAMEBUFFER, fbo)
	defer C.glBindFramebuffer(C.GL_FRAMEBUFFER, 0)
	C.glFramebufferTexture2D(C.GL_FRAMEBUFFER, C.GL_COLOR_ATTACHMENT0, C.GL_TEXTURE_2D, tex, 0)
	if st := C.glCheckFramebufferStatus(C.GL_FRAMEBUFFER); st != C.GL_FRAMEBUFFER_COMPLETE {
		return fmt.Errorf("screenshot: framebuffer incomplete (%#x)", st)
	}
	drawGL()
	if err := ctx.Frame(ops, gpu.OpenGLRenderTarget{V: uint(fbo)}, size); err != nil {
		return fmt.Errorf("screenshot: %w", err)
	}
	r := image.Rectangle{Max: size}
	ss := image.NewRGBA(r)
	C.glReadPixels(C.GLint(r.Min.X), C.GLint(r.Min.Y), C.GLint(r.Dx()), C.GLint(r.Dy()), C.GL_RGBA, C.GL_UNSIGNED_BYTE, unsafe.Pointer(&ss.Pix[0]))
	var buf bytes.Buffer
	if err := png.Encode(&buf, ss); err != nil {
		return fmt.Errorf("screenshot: %w", err)
	}
	const file = "screenshot.png"
	if err := os.WriteFile(file, buf.Bytes(), 0o644); err != nil {
		return fmt.Errorf("screenshot: %w", err)
	}
	fmt.Printf("wrote %q\n", file)
	return nil
}

func drawGL() {
	C.glClearColor(.5, .5, 0, 1)
	C.glClear(C.GL_COLOR_BUFFER_BIT | C.GL_DEPTH_BUFFER_BIT)
}

func drawUI(th *material.Theme, gtx layout.Context) layout.Dimensions {
	return layout.Center.Layout(gtx,
		material.Button(th, &btnScreenshot, "Screenshot").Layout,
	)
}

func createContext(ve app.ViewEvent, size image.Point) (*eglContext, error) {
	view, cleanup := nativeViewFor(ve, size)
	var nilv C.EGLNativeWindowType
	if view == nilv {
		return nil, fmt.Errorf("failed creating native view")
	}
	disp := getDisplay(ve)
	if disp == 0 {
		return nil, fmt.Errorf("eglGetPlatformDisplay failed: 0x%x", C.eglGetError())
	}
	var major, minor C.EGLint
	if ok := C.eglInitialize(disp, &major, &minor); ok != C.EGL_TRUE {
		return nil, fmt.Errorf("eglInitialize failed: 0x%x", C.eglGetError())
	}
	exts := strings.Split(C.GoString(C.eglQueryString(disp, C.EGL_EXTENSIONS)), " ")
	srgb := hasExtension(exts, "EGL_KHR_gl_colorspace")
	attribs := []C.EGLint{
		C.EGL_RENDERABLE_TYPE, C.EGL_OPENGL_ES2_BIT,
		C.EGL_SURFACE_TYPE, C.EGL_WINDOW_BIT,
		C.EGL_BLUE_SIZE, 8,
		C.EGL_GREEN_SIZE, 8,
		C.EGL_RED_SIZE, 8,
		C.EGL_CONFIG_CAVEAT, C.EGL_NONE,
	}
	if srgb {
		// Some drivers need alpha for sRGB framebuffers to work.
		attribs = append(attribs, C.EGL_ALPHA_SIZE, 8)
	}
	attribs = append(attribs, C.EGL_NONE)
	var (
		cfg     C.EGLConfig
		numCfgs C.EGLint
	)
	if ok := C.eglChooseConfig(disp, &attribs[0], &cfg, 1, &numCfgs); ok != C.EGL_TRUE {
		return nil, fmt.Errorf("eglChooseConfig failed: 0x%x", C.eglGetError())
	}
	if numCfgs == 0 {
		supportsNoCfg := hasExtension(exts, "EGL_KHR_no_config_context")
		if !supportsNoCfg {
			return nil, errors.New("eglChooseConfig returned no configs")
		}
	}
	ctxAttribs := []C.EGLint{
		C.EGL_CONTEXT_CLIENT_VERSION, 3,
		C.EGL_NONE,
	}
	ctx := C.eglCreateContext(disp, cfg, nil, &ctxAttribs[0])
	if ctx == nil {
		return nil, fmt.Errorf("eglCreateContext failed: 0x%x", C.eglGetError())
	}
	var surfAttribs []C.EGLint
	if srgb {
		surfAttribs = append(surfAttribs, C.EGL_GL_COLORSPACE, C.EGL_GL_COLORSPACE_SRGB)
	}
	surfAttribs = append(surfAttribs, C.EGL_NONE)
	surf := C.eglCreateWindowSurface(disp, cfg, view, &surfAttribs[0])
	if surf == nil {
		return nil, fmt.Errorf("eglCreateWindowSurface failed (0x%x)", C.eglGetError())
	}
	e := &eglContext{
		view:    view,
		disp:    disp,
		ctx:     ctx,
		surf:    surf,
		cleanup: cleanup,
	}
	return e, nil
}

func (c *eglContext) resizeSurface(ve app.ViewEvent, size image.Point) {
	nativeViewResize(ve, c.view, size)
}

func (c *eglContext) Release() {
	if c.ctx != nil {
		C.eglDestroyContext(c.disp, c.ctx)
	}
	if c.surf != nil {
		C.eglDestroySurface(c.disp, c.surf)
	}
	if c.cleanup != nil {
		c.cleanup()
	}
	*c = eglContext{}
}

func hasExtension(exts []string, ext string) bool {
	for _, e := range exts {
		if ext == e {
			return true
		}
	}
	return false
}
