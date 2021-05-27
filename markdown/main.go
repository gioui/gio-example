// SPDX-License-Identifier: Unlicense OR MIT

package main

// A simple Gio program. See https://gioui.org for more information.

import (
	"bytes"
	"image"
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"gioui.org/x/richtext"

	"gioui.org/font/gofont"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type GioRenderer struct {
	richtext.TextObjects

	Current richtext.TextObject
	Theme   *material.Theme
}

func NewRenderer(theme *material.Theme) renderer.NodeRenderer {
	l := material.Body1(theme, "")
	g := &GioRenderer{
		Theme: theme,
	}
	g.UpdateCurrent(l)
	return g
}

func (g *GioRenderer) CommitCurrent() {
	g.TextObjects = append(g.TextObjects, g.Current.DeepCopy())
}

func (g *GioRenderer) UpdateCurrent(l material.LabelStyle) {
	g.Current.Font = l.Font
	g.Current.Color = l.Color
	g.Current.Size = l.TextSize
}

func (g *GioRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	// blocks
	//
	reg.Register(ast.KindDocument, g.renderDocument)
	reg.Register(ast.KindHeading, g.renderHeading)
	reg.Register(ast.KindBlockquote, g.renderBlockquote)
	reg.Register(ast.KindCodeBlock, g.renderCodeBlock)
	reg.Register(ast.KindFencedCodeBlock, g.renderFencedCodeBlock)
	reg.Register(ast.KindHTMLBlock, g.renderHTMLBlock)
	reg.Register(ast.KindList, g.renderList)
	reg.Register(ast.KindListItem, g.renderListItem)
	reg.Register(ast.KindParagraph, g.renderParagraph)
	reg.Register(ast.KindTextBlock, g.renderTextBlock)
	reg.Register(ast.KindThematicBreak, g.renderThematicBreak)
	//
	//	// inlines
	//
	reg.Register(ast.KindAutoLink, g.renderAutoLink)
	reg.Register(ast.KindCodeSpan, g.renderCodeSpan)
	reg.Register(ast.KindEmphasis, g.renderEmphasis)
	reg.Register(ast.KindImage, g.renderImage)
	reg.Register(ast.KindLink, g.renderLink)
	reg.Register(ast.KindRawHTML, g.renderRawHTML)
	reg.Register(ast.KindText, g.renderText)
	reg.Register(ast.KindString, g.renderString)
}

func (g *GioRenderer) renderDocument(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	log.Println("renderDocument")
	return ast.WalkContinue, nil
}

func (g *GioRenderer) renderHeading(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	log.Println("renderHeading")
	n := node.(*ast.Heading)
	if entering {
		var l material.LabelStyle
		switch n.Level {
		case 1:
			l = material.H1(g.Theme, "")
		case 2:
			l = material.H2(g.Theme, "")
		case 3:
			l = material.H3(g.Theme, "")
		case 4:
			l = material.H4(g.Theme, "")
		case 5:
			l = material.H5(g.Theme, "")
		case 6:
			l = material.H6(g.Theme, "")
		}
		g.UpdateCurrent(l)
	} else {
	}
	return ast.WalkContinue, nil
}

func (g *GioRenderer) renderBlockquote(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	log.Println("renderBlockquote")
	return ast.WalkContinue, nil
}

func (g *GioRenderer) renderCodeBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	log.Println("renderCodeBlock")
	return ast.WalkContinue, nil
}

func (g *GioRenderer) renderFencedCodeBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	log.Println("renderFencedCodeBlock")
	return ast.WalkContinue, nil
}

func (g *GioRenderer) renderHTMLBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	log.Println("renderHTMLBlock")
	return ast.WalkContinue, nil
}

func (g *GioRenderer) renderList(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	log.Println("renderList")
	return ast.WalkContinue, nil
}
func (g *GioRenderer) renderListItem(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	log.Println("renderListItem")
	return ast.WalkContinue, nil
}
func (g *GioRenderer) renderParagraph(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	log.Println("renderParagraph")
	return ast.WalkContinue, nil
}
func (g *GioRenderer) renderTextBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	log.Println("renderTextBlock")
	return ast.WalkContinue, nil
}
func (g *GioRenderer) renderThematicBreak(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	log.Println("renderThematicBreak")
	return ast.WalkContinue, nil
}
func (g *GioRenderer) renderAutoLink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	log.Println("renderAutoLink")
	return ast.WalkContinue, nil
}
func (g *GioRenderer) renderCodeSpan(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	log.Println("renderCodeSpan")
	if entering {
		g.Current.Font.Variant = "Mono"
	} else {
		g.Current.Font.Variant = ""
	}
	return ast.WalkContinue, nil
}
func (g *GioRenderer) renderEmphasis(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	log.Println("renderEmphasis")
	n := node.(*ast.Emphasis)

	if entering {
		if n.Level == 2 {
			g.Current.Font.Weight = text.Bold
		} else {
			g.Current.Font.Style = text.Italic
		}
	} else {
		g.Current.Font.Style = text.Regular
		g.Current.Font.Weight = text.Normal
	}
	return ast.WalkContinue, nil
}
func (g *GioRenderer) renderImage(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	log.Println("renderImage")
	return ast.WalkContinue, nil
}

const urlMetadataKey = "url"

func (g *GioRenderer) renderLink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	log.Println("renderLink")
	n := node.(*ast.Link)
	if entering {
		g.Current.Color = g.Theme.ContrastBg
		g.Current.Clickable = true
		g.Current.SetMetadata("url", string(n.Destination))
	} else {
		g.Current.Color = g.Theme.Fg
		g.Current.Clickable = false
		g.Current.SetMetadata("url", "")
	}
	return ast.WalkContinue, nil
}
func (g *GioRenderer) renderRawHTML(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	log.Println("renderRawHTML")
	return ast.WalkContinue, nil
}
func (g *GioRenderer) renderText(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	log.Println("renderText")
	if !entering {
		return ast.WalkContinue, nil
	}
	n := node.(*ast.Text)
	segment := n.Segment
	content := segment.Value(source)
	g.Current.Content = string(content)
	g.CommitCurrent()

	return ast.WalkContinue, nil
}
func (g *GioRenderer) renderString(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	log.Println("renderString")
	if !entering {
		return ast.WalkContinue, nil
	}
	n := node.(*ast.String)
	g.Current.Content = string(n.Value)
	g.CommitCurrent()
	return ast.WalkContinue, nil
}

func (g *GioRenderer) Result() richtext.TextObjects {
	o := g.TextObjects
	g.TextObjects = nil
	return o
}
func main() {
	go func() {
		w := app.NewWindow()
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

type (
	C = layout.Context
	D = layout.Dimensions
)

func loop(w *app.Window) error {
	fontCollection := gofont.Collection()
	shaper := text.NewCache(fontCollection)
	th := material.NewTheme(fontCollection)
	nr := NewRenderer(th)
	var ops op.Ops

	md := goldmark.New(
		goldmark.WithRenderer(
			renderer.NewRenderer(
				renderer.WithNodeRenderers(
					util.PrioritizedValue{Value: nr, Priority: 0},
				),
			),
		),
	)
	var buf bytes.Buffer
	var ed widget.Editor
	var rs component.Resize
	var rendered richtext.TextObjects
	rs.Ratio = .5
	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			for o := rendered.Clicked(); o != nil; o = rendered.Clicked() {
				log.Println(o)
				if url := o.GetMetadata(urlMetadataKey); url != "" {
					log.Println("URL", url)
				}
			}

			for range ed.Events() {
				if err := md.Convert([]byte(ed.Text()), &buf); err != nil {
					panic(err)
				}
				rendered = nr.(*GioRenderer).Result()
			}
			rs.Layout(gtx,
				material.Editor(th, &ed, "markdown").Layout,
				func(gtx C) D {
					return rendered.Layout(gtx, shaper)
				},
				func(gtx C) D {
					rect := image.Rectangle{
						Max: image.Point{
							X: (gtx.Px(unit.Dp(4))),
							Y: (gtx.Constraints.Max.Y),
						},
					}
					paint.FillShape(gtx.Ops, color.NRGBA{A: 200}, clip.Rect(rect).Op())
					return D{Size: rect.Max}
				},
			)
			e.Frame(gtx.Ops)
		}
	}
}
