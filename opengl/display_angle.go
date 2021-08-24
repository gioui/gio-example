// SPDX-License-Identifier: Unlicense OR MIT

//go:build windows || darwin
// +build windows darwin

package main

import (
	"strings"

	"gioui.org/app"
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

func getDisplay(_ app.ViewEvent) C.EGLDisplay {
	var EGL_NO_DISPLAY C.EGLDisplay
	platformExts := strings.Split(C.GoString(C.eglQueryString(EGL_NO_DISPLAY, C.EGL_EXTENSIONS)), " ")
	platformType := C.EGLint(C.EGL_PLATFORM_ANGLE_TYPE_DEFAULT_ANGLE)
	if hasExtension(platformExts, "EGL_ANGLE_platform_angle_metal") {
		// The Metal backend works better than the OpenGL backend.
		platformType = C.EGL_PLATFORM_ANGLE_TYPE_METAL_ANGLE
	}
	attrs := []C.EGLint{
		C.EGL_PLATFORM_ANGLE_TYPE_ANGLE,
		platformType,
		C.EGL_NONE,
	}
	return C.eglGetPlatformDisplayEXT(C.EGL_PLATFORM_ANGLE_ANGLE, nil, &attrs[0])
}
