//go:build darwin
// +build darwin

//go:generate mkdir -p example.app/Contents/MacOS
//go:generate go build -o example.app/Contents/MacOS/example
//go:generate codesign -s - example.app

package main
