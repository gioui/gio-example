// Code generated by 'yaegi extract gioui.org/x/scroll'. DO NOT EDIT.

package giopkgs

import (
	"gioui.org/x/scroll"
	"go/constant"
	"go/token"
	"reflect"
)

func init() {
	Symbols["gioui.org/x/scroll"] = map[string]reflect.Value{
		// function, constant and variable definitions
		"DefaultBar": reflect.ValueOf(scroll.DefaultBar),
		"Horizontal": reflect.ValueOf(constant.MakeFromLiteral("1", token.INT, 0)),
		"Vertical":   reflect.ValueOf(constant.MakeFromLiteral("0", token.INT, 0)),

		// type definitions
		"Axis":       reflect.ValueOf((*scroll.Axis)(nil)),
		"Bar":        reflect.ValueOf((*scroll.Bar)(nil)),
		"C":          reflect.ValueOf((*scroll.C)(nil)),
		"D":          reflect.ValueOf((*scroll.D)(nil)),
		"Scrollable": reflect.ValueOf((*scroll.Scrollable)(nil)),
	}
}
