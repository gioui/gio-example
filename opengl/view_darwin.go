// SPDX-License-Identifier: Unlicense OR MIT

package main

import (
	"gioui.org/app"
)

/*
#include <EGL/egl.h>
*/
import "C"

func nativeViewFor(e app.ViewEvent) C.EGLNativeWindowType {
	return C.EGLNativeWindowType(e.Layer)
}
