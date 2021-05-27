// SPDX-License-Identifier: Unlicense OR MIT

package main

// A simple Gio program. See https://gioui.org for more information.

import (
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
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
	"github.com/inkeliz/giohyperlink"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type GioNodeRenderer struct {
	richtext.TextObjects

	Current      richtext.TextObject
	Theme        *material.Theme
	OrderedList  bool
	OrderedIndex int
}

func NewNodeRenderer(theme *material.Theme) *GioNodeRenderer {
	l := material.Body1(theme, "")
	g := &GioNodeRenderer{
		Theme: theme,
	}
	g.UpdateCurrent(l)
	return g
}

func (g *GioNodeRenderer) CommitCurrent() {
	g.TextObjects = append(g.TextObjects, g.Current.DeepCopy())
}

func (g *GioNodeRenderer) UpdateCurrent(l material.LabelStyle) {
	g.Current.Font = l.Font
	g.Current.Color = l.Color
	g.Current.Size = l.TextSize
}

func (g *GioNodeRenderer) AppendNewline() {
	if len(g.TextObjects) < 1 {
		return
	}
	g.TextObjects[len(g.TextObjects)-1].Content += "\n"
}

func (g *GioNodeRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
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

func (g *GioNodeRenderer) renderDocument(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	log.Println("renderDocument")
	return ast.WalkContinue, nil
}

func (g *GioNodeRenderer) renderHeading(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
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
		l := material.Body1(g.Theme, "")
		g.UpdateCurrent(l)
		g.AppendNewline()
	}
	return ast.WalkContinue, nil
}

func (g *GioNodeRenderer) renderBlockquote(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	log.Println("renderBlockquote")
	return ast.WalkContinue, nil
}

func (g *GioNodeRenderer) renderCodeBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	log.Println("renderCodeBlock")
	if entering {
		g.Current.Font.Variant = "Mono"
	} else {
		g.Current.Font.Variant = ""
	}
	return ast.WalkContinue, nil
}

func (g *GioNodeRenderer) renderFencedCodeBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	log.Println("renderFencedCodeBlock")
	n := node.(*ast.FencedCodeBlock)
	if entering {
		g.Current.Font.Variant = "Mono"
		lines := n.Lines()
		for i := 0; i < lines.Len(); i++ {
			line := lines.At(i)
			g.Current.Content = string(line.Value(source))
			g.CommitCurrent()
		}
	} else {
		g.Current.Font.Variant = ""
		g.AppendNewline()
	}
	return ast.WalkContinue, nil
}

func (g *GioNodeRenderer) renderHTMLBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	log.Println("renderHTMLBlock")
	if entering {
		g.Current.Font.Variant = "Mono"
	} else {
		g.Current.Font.Variant = ""
	}
	return ast.WalkContinue, nil
}

func (g *GioNodeRenderer) renderList(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	log.Println("renderList")
	n := node.(*ast.List)
	if entering {
		g.OrderedList = n.IsOrdered()
		g.OrderedIndex = 1
	} else {
		g.AppendNewline()
	}
	return ast.WalkContinue, nil
}

func (g *GioNodeRenderer) renderListItem(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	log.Println("renderListItem")
	if entering {
		if g.OrderedList {
			g.Current.Content = fmt.Sprintf(" %d. ", g.OrderedIndex)
			g.OrderedIndex++
		} else {
			g.Current.Content = " • "
		}
		g.CommitCurrent()
	} else if len(g.TextObjects) > 0 {
		g.AppendNewline()
	}

	return ast.WalkContinue, nil
}
func (g *GioNodeRenderer) renderParagraph(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	log.Println("renderParagraph")
	if !entering {
		g.AppendNewline()
		g.AppendNewline()
	}
	return ast.WalkContinue, nil
}
func (g *GioNodeRenderer) renderTextBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	log.Println("renderTextBlock")
	return ast.WalkContinue, nil
}
func (g *GioNodeRenderer) renderThematicBreak(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	log.Println("renderThematicBreak")
	return ast.WalkContinue, nil
}
func (g *GioNodeRenderer) renderAutoLink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	log.Println("renderAutoLink")
	n := node.(*ast.AutoLink)
	if entering {
		url := string(n.URL(source))
		g.Current.SetMetadata(urlMetadataKey, url)
		g.Current.Color = g.Theme.ContrastBg
		g.Current.Content = url
		g.CommitCurrent()
	} else {
		g.Current.SetMetadata(urlMetadataKey, "")
		g.Current.Color = g.Theme.Fg
	}
	return ast.WalkContinue, nil
}
func (g *GioNodeRenderer) renderCodeSpan(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	log.Println("renderCodeSpan")
	if entering {
		g.Current.Font.Variant = "Mono"
	} else {
		g.Current.Font.Variant = ""
	}
	return ast.WalkContinue, nil
}
func (g *GioNodeRenderer) renderEmphasis(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
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
func (g *GioNodeRenderer) renderImage(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	log.Println("renderImage")
	return ast.WalkContinue, nil
}

const urlMetadataKey = "url"

func (g *GioNodeRenderer) renderLink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
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
func (g *GioNodeRenderer) renderRawHTML(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	log.Println("renderRawHTML")
	return ast.WalkContinue, nil
}
func (g *GioNodeRenderer) renderText(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
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
func (g *GioNodeRenderer) renderString(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	log.Println("renderString")
	if !entering {
		return ast.WalkContinue, nil
	}
	n := node.(*ast.String)
	g.Current.Content = string(n.Value)
	g.CommitCurrent()
	return ast.WalkContinue, nil
}

func (g *GioNodeRenderer) Result() richtext.TextObjects {
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

type Renderer struct {
	md goldmark.Markdown
	nr *GioNodeRenderer
}

func NewRenderer(th *material.Theme) *Renderer {
	nr := NewNodeRenderer(th)
	md := goldmark.New(
		goldmark.WithRenderer(
			renderer.NewRenderer(
				renderer.WithNodeRenderers(
					util.PrioritizedValue{Value: nr, Priority: 0},
				),
			),
		),
	)
	return &Renderer{md: md, nr: nr}
}

func (r *Renderer) Render(src []byte) (richtext.TextObjects, error) {
	if err := r.md.Convert(src, ioutil.Discard); err != nil {
		return nil, err
	}
	return r.nr.Result(), nil
}

func loop(w *app.Window) error {
	fontCollection := gofont.Collection()
	shaper := text.NewCache(fontCollection)
	th := material.NewTheme(fontCollection)
	renderer := NewRenderer(th)
	var ops op.Ops

	var ed widget.Editor
	var rs component.Resize
	rs.Ratio = .5
	var rendered richtext.TextObjects
	inset := layout.UniformInset(unit.Dp(4))
	for {
		e := <-w.Events()
		giohyperlink.ListenEvents(e)
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			for o := rendered.Clicked(); o != nil; o = rendered.Clicked() {
				log.Println(o)
				if url := o.GetMetadata(urlMetadataKey); url != "" {
					giohyperlink.Open(url)
				}
			}

			for _, edEvent := range ed.Events() {
				if _, ok := edEvent.(widget.ChangeEvent); ok {
					rendered, _ = renderer.Render([]byte(ed.Text()))
				}
			}

			rs.Layout(gtx,
				func(gtx C) D { return inset.Layout(gtx, material.Editor(th, &ed, "markdown").Layout) },
				func(gtx C) D {
					return inset.Layout(gtx, func(gtx C) D {
						return rendered.Layout(gtx, shaper)
					})
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
