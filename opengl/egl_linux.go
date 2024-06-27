// SPDX-License-Identifier: Unlicense OR MIT

package main

import (
	"image"
	"unsafe"

	"gioui.org/app"
)

/*
#cgo linux pkg-config: egl wayland-egl
#cgo CFLAGS: -DEGL_NO_X11
#cgo LDFLAGS: -lEGL -lGLESv2

#include <EGL/egl.h>
#include <wayland-client.h>
#include <wayland-egl.h>
#include <GLES3/gl3.h>
#define EGL_EGLEXT_PROTOTYPES
#include <EGL/eglext.h>

*/
import "C"

func getDisplay(ve app.ViewEvent) C.EGLDisplay {
	switch ve := ve.(type) {
	case app.X11ViewEvent:
		return C.eglGetDisplay(C.EGLNativeDisplayType(ve.Display))
	case app.WaylandViewEvent:
		return C.eglGetDisplay(C.EGLNativeDisplayType(ve.Display))
	}
	panic("no display available")
}

func nativeViewFor(e app.ViewEvent, size image.Point) (C.EGLNativeWindowType, func()) {
	switch e := e.(type) {
	case app.X11ViewEvent:
		return C.EGLNativeWindowType(uintptr(e.Window)), func() {}
	case app.WaylandViewEvent:
		eglWin := C.wl_egl_window_create((*C.struct_wl_surface)(e.Surface), C.int(size.X), C.int(size.Y))
		return C.EGLNativeWindowType(uintptr(unsafe.Pointer(eglWin))), func() {
			C.wl_egl_window_destroy(eglWin)
		}
	}
	panic("no native view available")
}

func nativeViewResize(e app.ViewEvent, view C.EGLNativeWindowType, newSize image.Point) {
	switch e.(type) {
	case app.X11ViewEvent:
	case app.WaylandViewEvent:
		C.wl_egl_window_resize(*(**C.struct_wl_egl_window)(unsafe.Pointer(&view)), C.int(newSize.X), C.int(newSize.Y), 0, 0)
	}
}
