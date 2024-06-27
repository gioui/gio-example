// SPDX-License-Identifier: Unlicense OR MIT

package main

import (
	"image"
	"unsafe"

	"gioui.org/app"
)

/*
#include <EGL/egl.h>
*/
import "C"

func nativeViewFor(e app.ViewEvent, _ image.Point) (C.EGLNativeWindowType, func()) {
	return C.EGLNativeWindowType(unsafe.Pointer(e.(app.Win32ViewEvent).HWND)), func() {}
}
func nativeViewResize(e app.ViewEvent, view C.EGLNativeWindowType, newSize image.Point) {}
