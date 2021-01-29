package giopkgs

//go:generate go run github.com/traefik/yaegi/cmd/yaegi extract gioui.org/app
//go:generate go run github.com/traefik/yaegi/cmd/yaegi extract gioui.org/app/headless
//go:generate go run github.com/traefik/yaegi/cmd/yaegi extract gioui.org/f32
//go:generate go run github.com/traefik/yaegi/cmd/yaegi extract gioui.org/font/gofont
//go:generate go run github.com/traefik/yaegi/cmd/yaegi extract gioui.org/font/opentype
//go:generate go run github.com/traefik/yaegi/cmd/yaegi extract gioui.org/gesture
//go:generate go run github.com/traefik/yaegi/cmd/yaegi extract gioui.org/gpu
//go:generate go run github.com/traefik/yaegi/cmd/yaegi extract gioui.org/gpu/backend
//go:generate go run github.com/traefik/yaegi/cmd/yaegi extract gioui.org/gpu/gl
//go:generate go run github.com/traefik/yaegi/cmd/yaegi extract gioui.org/io/clipboard
//go:generate go run github.com/traefik/yaegi/cmd/yaegi extract gioui.org/io/event
//go:generate go run github.com/traefik/yaegi/cmd/yaegi extract gioui.org/io/key
//go:generate go run github.com/traefik/yaegi/cmd/yaegi extract gioui.org/io/pointer
//go:generate go run github.com/traefik/yaegi/cmd/yaegi extract gioui.org/io/profile
//go:generate go run github.com/traefik/yaegi/cmd/yaegi extract gioui.org/io/router
//go:generate go run github.com/traefik/yaegi/cmd/yaegi extract gioui.org/io/system
//go:generate go run github.com/traefik/yaegi/cmd/yaegi extract gioui.org/layout
//go:generate go run github.com/traefik/yaegi/cmd/yaegi extract gioui.org/op
//go:generate go run github.com/traefik/yaegi/cmd/yaegi extract gioui.org/op/clip
//go:generate go run github.com/traefik/yaegi/cmd/yaegi extract gioui.org/op/paint
//go:generate go run github.com/traefik/yaegi/cmd/yaegi extract gioui.org/text
//go:generate go run github.com/traefik/yaegi/cmd/yaegi extract gioui.org/unit
//go:generate go run github.com/traefik/yaegi/cmd/yaegi extract gioui.org/widget
//go:generate go run github.com/traefik/yaegi/cmd/yaegi extract gioui.org/widget/material
//go:generate go run github.com/traefik/yaegi/cmd/yaegi extract gioui.org/x/colorpicker
//go:generate go run github.com/traefik/yaegi/cmd/yaegi extract gioui.org/x/component
//go:generate go run github.com/traefik/yaegi/cmd/yaegi extract gioui.org/x/eventx
//go:generate go run github.com/traefik/yaegi/cmd/yaegi extract gioui.org/x/haptic
//go:generate go run github.com/traefik/yaegi/cmd/yaegi extract gioui.org/x/notify
//go:generate go run github.com/traefik/yaegi/cmd/yaegi extract gioui.org/x/outlay
//go:generate go run github.com/traefik/yaegi/cmd/yaegi extract gioui.org/x/profiling
//go:generate go run github.com/traefik/yaegi/cmd/yaegi extract gioui.org/x/scroll

import "reflect"

var Symbols = make(map[string]map[string]reflect.Value)
