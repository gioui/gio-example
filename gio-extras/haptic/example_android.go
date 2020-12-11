package main

import (
	"gioui.org/app"
	"gioui.org/io/event"
)

// ProcessPlatformEvent handles platform-specific event processing. If it
// consumed the provided event, it returns true. In this case, no further
// event processing should occur.
func ProcessPlatformEvent(event event.Event) bool {
	if ve, ok := event.(app.ViewEvent); ok {
		buzzer.SetView(ve.View)
		return true
	}
	return false
}
