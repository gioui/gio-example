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
	return C.EGLNativeWindowType(unsafe.Pointer(e.HWND)), func() {}
}
