// SPDX-License-Identifier: Unlicense OR MIT

package main

// A Gio program that demonstrates Gio widgets. See https://gioui.org for more information.

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/gesture"
	"gioui.org/gpu/headless"
	"gioui.org/io/pointer"
	"gioui.org/io/router"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/esimov/triangle"

	"golang.org/x/exp/shiny/materialdesign/icons"
)

var (
	screenshot = flag.String("screenshot", "", "save a screenshot to a file and exit")
	disable    = flag.Bool("disable", false, "disable all widgets")
)

type iconAndTextButton struct {
	theme  *material.Theme
	button *widget.Clickable
	icon   *widget.Icon
	word   string
}

func main() {
	flag.Parse()
	editor.SetText(longText)
	ic, err := widget.NewIcon(icons.ContentAdd)
	if err != nil {
		log.Fatal(err)
	}
	icon = ic
	progressIncrementer = make(chan float32)
	if *screenshot != "" {
		if err := saveScreenshot(*screenshot); err != nil {
			fmt.Fprintf(os.Stderr, "failed to save screenshot: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	go func() {
		for {
			time.Sleep(time.Second)
			progressIncrementer <- 0.1
		}
	}()

	go func() {
		w := app.NewWindow(app.Size(unit.Dp(800), unit.Dp(700)))
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func saveScreenshot(f string) error {
	const scale = 1.5
	sz := image.Point{X: 800 * scale, Y: 600 * scale}
	w, err := headless.NewWindow(sz.X, sz.Y)
	if err != nil {
		return err
	}
	gtx := layout.Context{
		Ops: new(op.Ops),
		Metric: unit.Metric{
			PxPerDp: scale,
			PxPerSp: scale,
		},
		Constraints: layout.Exact(sz),
		Queue:       new(router.Router),
	}
	th := material.NewTheme(gofont.Collection())
	kitchen(gtx, th)
	w.Frame(gtx.Ops)
	img, err := w.Screenshot()
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return err
	}
	return ioutil.WriteFile(f, buf.Bytes(), 0666)
}

func loop(w *app.Window) error {
	th := material.NewTheme(gofont.Collection())
	var shatter Shatter

	var ops op.Ops
	for {
		select {
		case e := <-w.Events():
			switch e := e.(type) {
			case system.DestroyEvent:
				return e.Err
			case system.FrameEvent:
				gtx := layout.NewContext(&ops, e)
				if *disable {
					gtx = gtx.Disabled()
				}
				if checkbox.Changed() {
					if checkbox.Value {
						transformTime = e.Now
					} else {
						transformTime = time.Time{}
					}
				}

				paint.Fill(gtx.Ops, color.NRGBA{A: 255})
				shatter.Layout(gtx, func(gtx C) D {
					return transformedKitchen(gtx, th)
				})
				e.Frame(gtx.Ops)
			}
		case p := <-progressIncrementer:
			progress += p
			if progress > 1 {
				progress = 0
			}
			w.Invalidate()
		}
	}
}

func transformedKitchen(gtx layout.Context, th *material.Theme) layout.Dimensions {
	if !transformTime.IsZero() {
		dt := float32(gtx.Now.Sub(transformTime).Seconds())
		angle := dt * .1
		op.InvalidateOp{}.Add(gtx.Ops)
		defer op.Save(gtx.Ops).Load()
		tr := f32.Affine2D{}
		tr = tr.Rotate(f32.Pt(300, 20), -angle)
		scale := 1.0 - dt*.5
		if scale < 0.5 {
			scale = 0.5
		}
		tr = tr.Scale(f32.Pt(300, 20), f32.Pt(scale, scale))
		offset := dt * 50
		if offset > 200 {
			offset = 200
		}
		tr = tr.Offset(f32.Pt(0, offset))
		op.Affine(tr).Add(gtx.Ops)
	}

	return kitchen(gtx, th)
}

var (
	editor     = new(widget.Editor)
	lineEditor = &widget.Editor{
		SingleLine: true,
		Submit:     true,
	}
	button            = new(widget.Clickable)
	greenButton       = new(widget.Clickable)
	iconTextButton    = new(widget.Clickable)
	iconButton        = new(widget.Clickable)
	flatBtn           = new(widget.Clickable)
	disableBtn        = new(widget.Clickable)
	radioButtonsGroup = new(widget.Enum)
	list              = &widget.List{
		List: layout.List{
			Axis: layout.Vertical,
		},
	}
	progress            = float32(0)
	progressIncrementer chan float32
	green               = true
	topLabel            = "Hello, Gio"
	icon                *widget.Icon
	checkbox            = new(widget.Bool)
	swtch               = new(widget.Bool)
	transformTime       time.Time
	float               = new(widget.Float)
)

type (
	D = layout.Dimensions
	C = layout.Context
)

func (b iconAndTextButton) Layout(gtx layout.Context) layout.Dimensions {
	return material.ButtonLayout(b.theme, b.button).Layout(gtx, func(gtx C) D {
		return layout.UniformInset(unit.Dp(12)).Layout(gtx, func(gtx C) D {
			iconAndLabel := layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}
			textIconSpacer := unit.Dp(5)

			layIcon := layout.Rigid(func(gtx C) D {
				return layout.Inset{Right: textIconSpacer}.Layout(gtx, func(gtx C) D {
					var d D
					if b.icon != nil {
						size := gtx.Px(unit.Dp(56)) - 2*gtx.Px(unit.Dp(16))
						gtx.Constraints = layout.Exact(image.Pt(size, size))
						d = b.icon.Layout(gtx, b.theme.ContrastFg)
					}
					return d
				})
			})

			layLabel := layout.Rigid(func(gtx C) D {
				return layout.Inset{Left: textIconSpacer}.Layout(gtx, func(gtx C) D {
					l := material.Body1(b.theme, b.word)
					l.Color = b.theme.Palette.ContrastFg
					return l.Layout(gtx)
				})
			})

			return iconAndLabel.Layout(gtx, layIcon, layLabel)
		})
	})
}

func kitchen(gtx layout.Context, th *material.Theme) layout.Dimensions {
	paint.FillShape(gtx.Ops, color.NRGBA{R: 255, G: 255, B: 255, A: 255}, clip.Rect{Max: gtx.Constraints.Max}.Op())
	for _, e := range lineEditor.Events() {
		if e, ok := e.(widget.SubmitEvent); ok {
			topLabel = e.Text
			lineEditor.SetText("")
		}
	}
	widgets := []layout.Widget{
		material.H3(th, topLabel).Layout,
		func(gtx C) D {
			gtx.Constraints.Max.Y = gtx.Px(unit.Dp(200))
			return material.Editor(th, editor, "Hint").Layout(gtx)
		},
		func(gtx C) D {
			e := material.Editor(th, lineEditor, "Hint")
			e.Font.Style = text.Italic
			border := widget.Border{Color: color.NRGBA{A: 0xff}, CornerRadius: unit.Dp(8), Width: unit.Px(2)}
			return border.Layout(gtx, func(gtx C) D {
				return layout.UniformInset(unit.Dp(8)).Layout(gtx, e.Layout)
			})
		},
		func(gtx C) D {
			gtx.Constraints.Min.Y = gtx.Px(unit.Dp(50))
			gtx.Constraints.Max.Y = gtx.Constraints.Min.Y

			dr := image.Rectangle{Max: gtx.Constraints.Min}
			defer op.Save(gtx.Ops).Load()
			paint.LinearGradientOp{
				Stop1:  layout.FPt(dr.Min),
				Stop2:  layout.FPt(dr.Max),
				Color1: color.NRGBA{R: 0x10, G: 0xff, B: 0x10, A: 0xFF},
				Color2: color.NRGBA{R: 0x10, G: 0x10, B: 0xff, A: 0xFF},
			}.Add(gtx.Ops)
			clip.Rect(dr).Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			return layout.Dimensions{
				Size: gtx.Constraints.Max,
			}
		},
		func(gtx C) D {
			in := layout.UniformInset(unit.Dp(8))
			return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return in.Layout(gtx, material.IconButton(th, iconButton, icon).Layout)
				}),
				layout.Rigid(func(gtx C) D {
					return in.Layout(gtx, iconAndTextButton{theme: th, icon: icon, word: "Icon", button: iconTextButton}.Layout)
				}),
				layout.Rigid(func(gtx C) D {
					return in.Layout(gtx, func(gtx C) D {
						for button.Clicked() {
							green = !green
						}
						dims := material.Button(th, button, "Click me!").Layout(gtx)
						pointer.CursorNameOp{Name: pointer.CursorPointer}.Add(gtx.Ops)
						return dims
					})
				}),
				layout.Rigid(func(gtx C) D {
					return in.Layout(gtx, func(gtx C) D {
						l := "Green"
						if !green {
							l = "Blue"
						}
						btn := material.Button(th, greenButton, l)
						if green {
							btn.Background = color.NRGBA{A: 0xff, R: 0x9e, G: 0x9d, B: 0x24}
						}
						return btn.Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return in.Layout(gtx, func(gtx C) D {
						return material.Clickable(gtx, flatBtn, func(gtx C) D {
							return layout.UniformInset(unit.Dp(12)).Layout(gtx, func(gtx C) D {
								flatBtnText := material.Body1(th, "Flat")
								if gtx.Queue == nil {
									flatBtnText.Color.A = 150
								}
								return layout.Center.Layout(gtx, flatBtnText.Layout)
							})
						})
					})
				}),
				layout.Rigid(material.ProgressCircle(th, progress).Layout),
			)
		},
		material.ProgressBar(th, progress).Layout,
		func(gtx C) D {
			return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(
					material.CheckBox(th, checkbox, "Transform").Layout,
				),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Left: unit.Dp(16)}.Layout(gtx,
						material.Switch(th, swtch).Layout,
					)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Left: unit.Dp(16)}.Layout(gtx, func(gtx C) D {
						text := "enabled"
						if !swtch.Value {
							text = "disabled"
							gtx = gtx.Disabled()
						}
						btn := material.Button(th, disableBtn, text)
						return btn.Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Left: unit.Dp(16)}.Layout(gtx, func(gtx C) D {
						if !swtch.Value {
							return D{}
						}
						return material.Loader(th).Layout(gtx)
					})
				}),
			)
		},
		func(gtx C) D {
			return layout.Flex{}.Layout(gtx,
				layout.Rigid(material.RadioButton(th, radioButtonsGroup, "r1", "RadioButton1").Layout),
				layout.Rigid(material.RadioButton(th, radioButtonsGroup, "r2", "RadioButton2").Layout),
				layout.Rigid(material.RadioButton(th, radioButtonsGroup, "r3", "RadioButton3").Layout),
			)
		},
		func(gtx C) D {
			return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
				layout.Flexed(1, material.Slider(th, float, 0, 2*math.Pi).Layout),
				layout.Rigid(func(gtx C) D {
					return layout.UniformInset(unit.Dp(8)).Layout(gtx,
						material.Body1(th, fmt.Sprintf("%.2f", float.Value)).Layout,
					)
				}),
			)
		},
	}

	return material.List(th, list).Layout(gtx, len(widgets), func(gtx C, i int) D {
		return layout.UniformInset(unit.Dp(16)).Layout(gtx, widgets[i])
	})
}

const longText = `1. I learned from my grandfather, Verus, to use good manners, and to
put restraint on anger. 2. In the famous memory of my father I had a
pattern of modesty and manliness. 3. Of my mother I learned to be
pious and generous; to keep myself not only from evil deeds, but even
from evil thoughts; and to live with a simplicity which is far from
customary among the rich. 4. I owe it to my great-grandfather that I
did not attend public lectures and discussions, but had good and able
teachers at home; and I owe him also the knowledge that for things of
this nature a man should count no expense too great.

5. My tutor taught me not to favour either green or blue at the
chariot races, nor, in the contests of gladiators, to be a supporter
either of light or heavy armed. He taught me also to endure labour;
not to need many things; to serve myself without troubling others; not
to intermeddle in the affairs of others, and not easily to listen to
slanders against them.

6. Of Diognetus I had the lesson not to busy myself about vain things;
not to credit the great professions of such as pretend to work
wonders, or of sorcerers about their charms, and their expelling of
Demons and the like; not to keep quails (for fighting or divination),
nor to run after such things; to suffer freedom of speech in others,
and to apply myself heartily to philosophy. Him also I must thank for
my hearing first Bacchius, then Tandasis and Marcianus; that I wrote
dialogues in my youth, and took a liking to the philosopher's pallet
and skins, and to the other things which, by the Grecian discipline,
belong to that profession.

7. To Rusticus I owe my first apprehensions that my nature needed
reform and cure; and that I did not fall into the ambition of the
common Sophists, either by composing speculative writings or by
declaiming harangues of exhortation in public; further, that I never
strove to be admired by ostentation of great patience in an ascetic
life, or by display of activity and application; that I gave over the
study of rhetoric, poetry, and the graces of language; and that I did
not pace my house in my senatorial robes, or practise any similar
affectation. I observed also the simplicity of style in his letters,
particularly in that which he wrote to my mother from Sinuessa. I
learned from him to be easily appeased, and to be readily reconciled
with those who had displeased me or given cause of offence, so soon as
they inclined to make their peace; to read with care; not to rest
satisfied with a slight and superficial knowledge; nor quickly to
assent to great talkers. I have him to thank that I met with the
discourses of Epictetus, which he furnished me from his own library.

8. From Apollonius I learned true liberty, and tenacity of purpose; to
regard nothing else, even in the smallest degree, but reason always;
and always to remain unaltered in the agonies of pain, in the losses
of children, or in long diseases. He afforded me a living example of
how the same man can, upon occasion, be most yielding and most
inflexible. He was patient in exposition; and, as might well be seen,
esteemed his fine skill and ability in teaching others the principles
of philosophy as the least of his endowments. It was from him that I
learned how to receive from friends what are thought favours without
seeming humbled by the giver or insensible to the gift.`

type Fragment struct {
	A, B, C      f32.Point
	Velocity     f32.Point
	AngV         float32
	Displacement f32.Point
	AngDisp      float32
}

func (f *Fragment) Layout(gtx C, call op.CallOp) {
	defer op.Save(gtx.Ops).Load()
	f.Displacement = f.Displacement.Add(f.Velocity)
	f.AngDisp += f.AngV
	for f.AngDisp > 2*math.Pi {
		f.AngDisp -= 2 * math.Pi
	}
	for f.AngDisp < 0 {
		f.AngDisp += 2 * math.Pi
	}
	var p clip.Path
	p.Begin(gtx.Ops)
	p.MoveTo(f.A)
	p.LineTo(f.B)
	p.LineTo(f.C)
	p.Close()
	o := clip.Outline{Path: p.End()}
	op.Affine(
		f32.Affine2D{}.Rotate(f.center(), f.AngDisp).Offset(f.Displacement),
	).Add(gtx.Ops)
	o.Op().Add(gtx.Ops)
	call.Add(gtx.Ops)
}

func max(x float32, xs ...float32) float32 {
	for _, o := range xs {
		if o > x {
			x = o
		}
	}
	return x
}

func min(x float32, xs ...float32) float32 {
	for _, o := range xs {
		if o < x {
			x = o
		}
	}
	return x
}

func mid(a, b f32.Point) f32.Point {
	return f32.Point{
		X: (max(a.X, b.X) - min(a.X, b.X)) * .5,
		Y: (max(a.Y, b.Y) - min(a.Y, b.Y)) * .5,
	}
}

func (f *Fragment) center() f32.Point {
	return f32.Point{
		X: (f.A.X + f.B.X + f.C.X) / 3,
		Y: (f.A.Y + f.B.Y + f.C.Y) / 3,
	}
	return f.A
}

// Shatter lays out another widget with a passthrough click area on top.
// The first click will cause the presentation of the provided widget
// to break into moving, spinning fragments of UI.
type Shatter struct {
	gesture.Click
	broken    bool
	origin    f32.Point
	fragments []Fragment
	triangle.Delaunay
}

func (s *Shatter) Layout(gtx C, w layout.Widget) D {
	macro := op.Record(gtx.Ops)
	constraints := gtx.Constraints
	dims := layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx C) D {
			gtx.Constraints = constraints
			return w(gtx)
		}),
		layout.Expanded(func(gtx C) D {
			if !s.broken {
				pointer.PassOp{Pass: true}.Add(gtx.Ops)
				pointer.Rect(image.Rectangle{Max: gtx.Constraints.Min}).Add(gtx.Ops)
				s.Click.Add(gtx.Ops)
			}
			return D{Size: gtx.Constraints.Min}
		}),
	)
	call := macro.Stop()

	for _, click := range s.Click.Events(gtx) {
		switch click.Type {
		case gesture.TypeClick:
			if !s.broken {
				s.origin = click.Position
				s.broken = true
				s.chooseFragments(gtx, click.Position)
			}
		}
	}

	if !s.broken {
		call.Add(gtx.Ops)
		return dims
	}
	op.InvalidateOp{}.Add(gtx.Ops)
	for i := range s.fragments {
		f := &s.fragments[i]
		f.Layout(gtx, call)
	}

	return dims
}

func randomPoints(n int, in image.Point) []image.Point {
	out := make([]image.Point, n)
	for i := range out {
		out[i] = image.Point{
			X: rand.Intn(in.X),
			Y: rand.Intn(in.Y),
		}
	}
	return out
}

func (s *Shatter) chooseFragments(gtx C, clickPos f32.Point) {
	s.fragments = s.fragments[:0]

	var points []image.Point
	points = append(points, image.Point{
		X: int(clickPos.X),
		Y: int(clickPos.Y),
	})

	points = append(points, randomPoints(10, gtx.Constraints.Max)...)

	triangles := s.Delaunay.Init(gtx.Constraints.Max.X, gtx.Constraints.Max.Y).Insert(points).GetTriangles()

	for _, t := range triangles {
		f := Fragment{
			A:    Pt(t.Nodes[0]),
			B:    Pt(t.Nodes[1]),
			C:    Pt(t.Nodes[2]),
			AngV: math.Pi / 500 * (.5 - rand.Float32()),
		}
		center := f.center()
		vec := center.Sub(clickPos)
		f.Velocity = vec.Mul(1 / Distance(f32.Point{}, vec)).Mul(3)
		s.fragments = append(s.fragments, f)
	}
}

func Pt(n triangle.Node) f32.Point {
	return f32.Pt(float32(n.X), float32(n.Y))
}

func Distance(a, b f32.Point) float32 {
	deltaX := float64(b.X - a.X)
	deltaY := float64(b.Y - a.Y)
	return float32(math.Sqrt(deltaX*deltaX + deltaY*deltaY))
}
