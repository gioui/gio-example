// SPDX-License-Identifier: Unlicense OR MIT

//go:build linux && nowayland
// +build linux,nowayland

package main

import "gioui.org/app"

/*
#cgo CFLAGS: -DEGL_NO_X11
#cgo LDFLAGS: -lEGL -lGLESv2

#include <EGL/egl.h>
#include <GLES3/gl3.h>
#define EGL_EGLEXT_PROTOTYPES
#include <EGL/eglext.h>

*/
import "C"

func getDisplay(ve app.ViewEvent) C.EGLDisplay {
	return C.eglGetDisplay(C.EGLNativeDisplayType(ve.Display))
}

func nativeViewFor(e app.ViewEvent) C.EGLNativeWindowType {
	return C.EGLNativeWindowType(uintptr(e.Window))
}
