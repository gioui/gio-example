//go:build !android
// +build !android

package main

import "gioui.org/io/event"

func ProcessPlatformEvent(event event.Event) bool {
	return false
}
