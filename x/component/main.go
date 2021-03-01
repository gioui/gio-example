package main

import (
	"flag"
	"image/color"
	"log"
	"os"
	"time"
	"unicode"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"golang.org/x/exp/shiny/materialdesign/icons"

	"gioui.org/x/component"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

var MenuIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.NavigationMenu)
	return icon
}()

var HomeIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.ActionHome)
	return icon
}()

var SettingsIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.ActionSettings)
	return icon
}()

var OtherIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.ActionHelp)
	return icon
}()

var HeartIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.ActionFavorite)
	return icon
}()

var PlusIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.ContentAdd)
	return icon
}()

var EditIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.ContentCreate)
	return icon
}()

var barOnBottom bool

func main() {
	flag.BoolVar(&barOnBottom, "bottom-bar", false, "place the app bar on the bottom of the screen instead of the top")
	flag.Parse()
	go func() {
		w := app.NewWindow()
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

const (
	settingNameColumnWidth    = .3
	settingDetailsColumnWidth = 1 - settingNameColumnWidth
)

func LayoutAppBarPage(gtx C) D {
	return layout.Flex{
		Alignment: layout.Middle,
		Axis:      layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return inset.Layout(gtx, material.Body1(th, `The app bar widget provides a consistent interface element for triggering navigation and page-specific actions.

The controls below allow you to see the various features available in our App Bar implementation.`).Layout)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Alignment: layout.Baseline}.Layout(gtx,
				layout.Flexed(settingNameColumnWidth, func(gtx C) D {
					return inset.Layout(gtx, material.Body1(th, "Contextual App Bar").Layout)
				}),
				layout.Flexed(settingDetailsColumnWidth, func(gtx C) D {
					if contextBtn.Clicked() {
						bar.SetContextualActions(
							[]component.AppBarAction{
								component.SimpleIconAction(th, &red, HeartIcon,
									component.OverflowAction{
										Name: "House",
										Tag:  &red,
									},
								),
							},
							[]component.OverflowAction{
								{
									Name: "foo",
									Tag:  &blue,
								},
								{
									Name: "bar",
									Tag:  &green,
								},
							},
						)
						bar.ToggleContextual(gtx.Now, "Contextual Title")
					}
					return material.Button(th, &contextBtn, "Trigger").Layout(gtx)
				}),
			)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
				layout.Flexed(settingNameColumnWidth, func(gtx C) D {
					return inset.Layout(gtx, material.Body1(th, "Bottom App Bar").Layout)
				}),
				layout.Flexed(settingDetailsColumnWidth, func(gtx C) D {
					if bottomBar.Changed() {
						if bottomBar.Value {
							nav.Anchor = component.Bottom
						} else {
							nav.Anchor = component.Top
						}
					}

					return inset.Layout(gtx, material.Switch(th, &bottomBar).Layout)
				}),
			)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
				layout.Flexed(settingNameColumnWidth, func(gtx C) D {
					return inset.Layout(gtx, material.Body1(th, "Custom Navigation Icon").Layout)
				}),
				layout.Flexed(settingDetailsColumnWidth, func(gtx C) D {
					if customNavIcon.Changed() {
						if customNavIcon.Value {
							bar.NavigationIcon = HomeIcon
						} else {
							bar.NavigationIcon = MenuIcon
						}
					}
					return inset.Layout(gtx, material.Switch(th, &customNavIcon).Layout)
				}),
			)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Alignment: layout.Baseline}.Layout(gtx,
				layout.Flexed(settingNameColumnWidth, func(gtx C) D {
					return inset.Layout(gtx, material.Body1(th, "Animated Resize").Layout)
				}),
				layout.Flexed(settingDetailsColumnWidth, func(gtx C) D {
					return inset.Layout(gtx, material.Body2(th, "Resize the width of your screen to see app bar actions collapse into or emerge from the overflow menu (as size permits).").Layout)
				}),
			)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Alignment: layout.Baseline}.Layout(gtx,
				layout.Flexed(settingNameColumnWidth, func(gtx C) D {
					return inset.Layout(gtx, material.Body1(th, "Custom Action Buttons").Layout)
				}),
				layout.Flexed(settingDetailsColumnWidth, func(gtx C) D {
					if heartBtn.Clicked() {
						favorited = !favorited
					}
					return inset.Layout(gtx, material.Body2(th, "Click the heart action to see custom button behavior.").Layout)
				}),
			)
		}),
	)
}

func LayoutNavDrawerPage(gtx C) D {
	return layout.Flex{
		Alignment: layout.Middle,
		Axis:      layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return inset.Layout(gtx, material.Body1(th, `The nav drawer widget provides a consistent interface element for navigation.

The controls below allow you to see the various features available in our Navigation Drawer implementation.`).Layout)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
				layout.Flexed(settingNameColumnWidth, func(gtx C) D {
					return inset.Layout(gtx, material.Body1(th, "Use non-modal drawer").Layout)
				}),
				layout.Flexed(settingDetailsColumnWidth, func(gtx C) D {
					if nonModalDrawer.Changed() {
						if nonModalDrawer.Value {
							navAnim.Appear(gtx.Now)
						} else {
							navAnim.Disappear(gtx.Now)
						}
					}
					return inset.Layout(gtx, material.Switch(th, &nonModalDrawer).Layout)
				}),
			)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Alignment: layout.Baseline}.Layout(gtx,
				layout.Flexed(settingNameColumnWidth, func(gtx C) D {
					return inset.Layout(gtx, material.Body1(th, "Drag to Close").Layout)
				}),
				layout.Flexed(settingDetailsColumnWidth, func(gtx C) D {
					return inset.Layout(gtx, material.Body2(th, "You can close the modal nav drawer by dragging it to the left.").Layout)
				}),
			)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Alignment: layout.Baseline}.Layout(gtx,
				layout.Flexed(settingNameColumnWidth, func(gtx C) D {
					return inset.Layout(gtx, material.Body1(th, "Touch Scrim to Close").Layout)
				}),
				layout.Flexed(settingDetailsColumnWidth, func(gtx C) D {
					return inset.Layout(gtx, material.Body2(th, "You can close the modal nav drawer touching anywhere in the translucent scrim to the right.").Layout)
				}),
			)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Alignment: layout.Baseline}.Layout(gtx,
				layout.Flexed(settingNameColumnWidth, func(gtx C) D {
					return inset.Layout(gtx, material.Body1(th, "Bottom content anchoring").Layout)
				}),
				layout.Flexed(settingDetailsColumnWidth, func(gtx C) D {
					return inset.Layout(gtx, material.Body2(th, "If you toggle support for the bottom app bar in the App Bar settings, nav drawer content will anchor to the bottom of the drawer area instead of the top.").Layout)
				}),
			)
		}),
	)
}

const (
	sponsorEliasURL          = "https://github.com/sponsors/eliasnaur"
	sponsorChrisURLGitHub    = "https://github.com/sponsors/whereswaldon"
	sponsorChrisURLLiberapay = "https://liberapay.com/whereswaldon/"
)

func LayoutAboutPage(gtx C) D {
	th := *th
	th.Palette = currentAccent
	return layout.Flex{
		Alignment: layout.Middle,
		Axis:      layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return inset.Layout(gtx, material.Body1(&th, `This library implements material design components from https://material.io using https://gioui.org.

If you like this library and work like it, please consider sponsoring Elias and/or Chris!`).Layout)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
				layout.Flexed(settingDetailsColumnWidth, func(gtx C) D {
					return inset.Layout(gtx, material.Body1(&th, "Try another theme:").Layout)
				}),
				layout.Flexed(settingNameColumnWidth, func(gtx C) D {
					return inset.Layout(gtx, material.Switch(&th, &alternatePalette).Layout)
				}),
			)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
				layout.Flexed(settingDetailsColumnWidth, func(gtx C) D {
					return inset.Layout(gtx, material.Body1(&th, "Elias Naur can be sponsored on GitHub at "+sponsorEliasURL).Layout)
				}),
				layout.Flexed(settingNameColumnWidth, func(gtx C) D {
					if eliasCopyButton.Clicked() {
						clipboardRequests <- sponsorEliasURL
					}
					return inset.Layout(gtx, material.Button(&th, &eliasCopyButton, "Copy Sponsorship URL").Layout)
				}),
			)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
				layout.Flexed(settingDetailsColumnWidth, func(gtx C) D {
					return inset.Layout(gtx, material.Body1(&th, "Chris Waldon can be sponsored on GitHub at "+sponsorChrisURLGitHub+" and on Liberapay at "+sponsorChrisURLLiberapay).Layout)
				}),
				layout.Flexed(settingNameColumnWidth, func(gtx C) D {
					if chrisCopyButtonGH.Clicked() {
						clipboardRequests <- sponsorChrisURLGitHub
					}
					if chrisCopyButtonLP.Clicked() {
						clipboardRequests <- sponsorChrisURLLiberapay
					}
					return inset.Layout(gtx, func(gtx C) D {
						return layout.Flex{}.Layout(gtx,
							layout.Flexed(.5, material.Button(&th, &chrisCopyButtonGH, "Copy GitHub URL").Layout),
							layout.Flexed(.5, material.Button(&th, &chrisCopyButtonLP, "Copy Liberapay URL").Layout),
						)
					})
				}),
			)
		}),
	)
}

func LayoutTextFieldPage(gtx C) D {
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(
		gtx,
		layout.Rigid(func(gtx C) D {
			nameInput.Alignment = inputAlignment
			return nameInput.Layout(gtx, th, "Name")
		}),
		layout.Rigid(func(gtx C) D {
			return inset.Layout(gtx, material.Body2(th, "Responds to hover events.").Layout)
		}),
		layout.Rigid(func(gtx C) D {
			addressInput.Alignment = inputAlignment
			return addressInput.Layout(gtx, th, "Address")
		}),
		layout.Rigid(func(gtx C) D {
			return inset.Layout(gtx, material.Body2(th, "Label animates properly when you click to select the text field.").Layout)
		}),
		layout.Rigid(func(gtx C) D {
			priceInput.Prefix = func(gtx C) D {
				th := *th
				th.Palette.Fg = color.NRGBA{R: 100, G: 100, B: 100, A: 255}
				return material.Label(&th, th.TextSize, "$").Layout(gtx)
			}
			priceInput.Suffix = func(gtx C) D {
				th := *th
				th.Palette.Fg = color.NRGBA{R: 100, G: 100, B: 100, A: 255}
				return material.Label(&th, th.TextSize, ".00").Layout(gtx)
			}
			priceInput.SingleLine = true
			priceInput.Alignment = inputAlignment
			return priceInput.Layout(gtx, th, "Price")
		}),
		layout.Rigid(func(gtx C) D {
			return inset.Layout(gtx, material.Body2(th, "Can have prefix and suffix elements.").Layout)
		}),
		layout.Rigid(func(gtx C) D {
			if err := func() string {
				for _, r := range numberInput.Text() {
					if !unicode.IsDigit(r) {
						return "Must contain only digits"
					}
				}
				return ""
			}(); err != "" {
				numberInput.SetError(err)
			} else {
				numberInput.ClearError()
			}
			numberInput.SingleLine = true
			numberInput.Alignment = inputAlignment
			return numberInput.Layout(gtx, th, "Number")
		}),
		layout.Rigid(func(gtx C) D {
			return inset.Layout(gtx, material.Body2(th, "Can be validated.").Layout)
		}),
		layout.Rigid(func(gtx C) D {
			if tweetInput.TextTooLong() {
				tweetInput.SetError("Too many characters")
			} else {
				tweetInput.ClearError()
			}
			tweetInput.CharLimit = 128
			tweetInput.Helper = "Tweets have a limited character count"
			tweetInput.Alignment = inputAlignment
			return tweetInput.Layout(gtx, th, "Tweet")
		}),
		layout.Rigid(func(gtx C) D {
			return inset.Layout(gtx, material.Body2(th, "Can have a character counter and help text.").Layout)
		}),
		layout.Rigid(func(gtx C) D {
			if inputAlignmentEnum.Changed() {
				switch inputAlignmentEnum.Value {
				case layout.Start.String():
					inputAlignment = layout.Start
				case layout.Middle.String():
					inputAlignment = layout.Middle
				case layout.End.String():
					inputAlignment = layout.End
				default:
					inputAlignment = layout.Start
				}
				op.InvalidateOp{}.Add(gtx.Ops)
			}
			return inset.Layout(
				gtx,
				func(gtx C) D {
					return layout.Flex{
						Axis: layout.Vertical,
					}.Layout(
						gtx,
						layout.Rigid(func(gtx C) D {
							return material.Body2(th, "Text Alignment").Layout(gtx)
						}),
						layout.Rigid(func(gtx C) D {
							return layout.Flex{
								Axis: layout.Vertical,
							}.Layout(
								gtx,
								layout.Rigid(func(gtx C) D {
									return material.RadioButton(
										th,
										&inputAlignmentEnum,
										layout.Start.String(),
										"Start",
									).Layout(gtx)
								}),
								layout.Rigid(func(gtx C) D {
									return material.RadioButton(
										th,
										&inputAlignmentEnum,
										layout.Middle.String(),
										"Middle",
									).Layout(gtx)
								}),
								layout.Rigid(func(gtx C) D {
									return material.RadioButton(
										th,
										&inputAlignmentEnum,
										layout.End.String(),
										"End",
									).Layout(gtx)
								}),
							)
						}),
					)
				},
			)
		}),
		layout.Rigid(func(gtx C) D {
			return inset.Layout(gtx, material.Body2(th, "This text field implementation was contributed by Jack Mordaunt. Thanks Jack!").Layout)
		}),
	)
}

type Page struct {
	layout func(layout.Context) layout.Dimensions
	component.NavItem
	Actions  []component.AppBarAction
	Overflow []component.OverflowAction

	// laying each page out within a layout.List enables scrolling for the page
	// content.
	layout.List
}

var (
	// initialize channel to send clipboard content requests on
	clipboardRequests = make(chan string, 1)

	// initialize modal layer to draw modal components
	modal   = component.NewModal()
	navAnim = component.VisibilityAnimation{
		Duration: time.Millisecond * 100,
		State:    component.Invisible,
	}
	nav      = component.NewNav("Navigation Drawer", "This is an example.")
	modalNav = component.ModalNavFrom(&nav, modal)

	bar = component.NewAppBar(modal)

	inset              = layout.UniformInset(unit.Dp(8))
	th                 = material.NewTheme(gofont.Collection())
	lightPalette       = th.Palette
	lightPaletteAccent = func() material.Palette {
		out := th.Palette
		out.ContrastBg = color.NRGBA{A: 0xff, R: 0x9e, G: 0x9d, B: 0x24}
		return out
	}()
	altPalette = func() material.Palette {
		out := th.Palette
		out.Bg = color.NRGBA{R: 0xff, G: 0xfb, B: 0xe6, A: 0xff}
		out.Fg = color.NRGBA{A: 0xff}
		out.ContrastBg = color.NRGBA{R: 0x35, G: 0x69, B: 0x59, A: 0xff}
		return out
	}()
	altPaletteAccent = func() material.Palette {
		out := th.Palette
		out.Bg = color.NRGBA{R: 0xff, G: 0xfb, B: 0xe6, A: 0xff}
		out.Fg = color.NRGBA{A: 0xff}
		out.ContrastBg = color.NRGBA{R: 0xfd, G: 0x55, B: 0x23, A: 0xff}
		out.ContrastFg = out.Fg
		return out
	}()
	currentAccent material.Palette = lightPaletteAccent

	heartBtn, plusBtn, exampleOverflowState               widget.Clickable
	red, green, blue                                      widget.Clickable
	contextBtn                                            widget.Clickable
	eliasCopyButton, chrisCopyButtonGH, chrisCopyButtonLP widget.Clickable
	bottomBar                                             widget.Bool
	customNavIcon                                         widget.Bool
	nonModalDrawer                                        widget.Bool
	alternatePalette                                      widget.Bool
	favorited                                             bool
	inputAlignment                                        layout.Alignment
	inputAlignmentEnum                                    widget.Enum
	nameInput                                             component.TextField
	addressInput                                          component.TextField
	priceInput                                            component.TextField
	tweetInput                                            component.TextField
	numberInput                                           component.TextField

	pages = []Page{
		{
			NavItem: component.NavItem{
				Name: "App Bar Features",
				Icon: HomeIcon,
			},
			layout: LayoutAppBarPage,
			Actions: []component.AppBarAction{
				{
					OverflowAction: component.OverflowAction{
						Name: "Favorite",
						Tag:  &heartBtn,
					},
					Layout: func(gtx layout.Context, bg, fg color.NRGBA) layout.Dimensions {
						btn := component.SimpleIconButton(th, &heartBtn, HeartIcon)
						btn.Background = bg
						if favorited {
							btn.Color = color.NRGBA{R: 200, A: 255}
						} else {
							btn.Color = fg
						}
						return btn.Layout(gtx)
					},
				},
				component.SimpleIconAction(th, &plusBtn, PlusIcon,
					component.OverflowAction{
						Name: "Create",
						Tag:  &plusBtn,
					},
				),
			},
			Overflow: []component.OverflowAction{
				{
					Name: "Example 1",
					Tag:  &exampleOverflowState,
				},
				{
					Name: "Example 2",
					Tag:  &exampleOverflowState,
				},
			},
		},
		{
			NavItem: component.NavItem{
				Name: "Nav Drawer Features",
				Icon: SettingsIcon,
			},
			layout: LayoutNavDrawerPage,
		},
		{
			NavItem: component.NavItem{
				Name: "Text Field Features",
				Icon: EditIcon,
			},
			layout: LayoutTextFieldPage,
		},
		{
			NavItem: component.NavItem{
				Name: "About this library",
				Icon: OtherIcon,
			},
			layout:  LayoutAboutPage,
			Actions: []component.AppBarAction{},
		},
	}
)

func loop(w *app.Window) error {
	var ops op.Ops

	bar.NavigationIcon = MenuIcon
	if barOnBottom {
		bar.Anchor = component.Bottom
		nav.Anchor = component.Bottom
	}

	// assign navigation tags and configure navigation bar with all pages
	for i := range pages {
		page := &pages[i]
		page.List.Axis = layout.Vertical
		page.NavItem.Tag = i
		nav.AddNavItem(page.NavItem)
	}

	// configure app bar initial state
	page := pages[nav.CurrentNavDestination().(int)]
	bar.Title = page.Name
	bar.SetActions(page.Actions, page.Overflow)

	for {
		select {
		case content := <-clipboardRequests:
			w.WriteClipboard(content)
		case e := <-w.Events():
			switch e := e.(type) {
			case system.DestroyEvent:
				return e.Err
			case system.FrameEvent:
				skipInset := e
				skipInset.Insets = system.Insets{}
				gtx := layout.NewContext(&ops, skipInset)
				for _, event := range bar.Events(gtx) {
					switch event := event.(type) {
					case component.AppBarNavigationClicked:
						if nonModalDrawer.Value {
							navAnim.ToggleVisibility(gtx.Now)
						} else {
							modalNav.Appear(gtx.Now)
							navAnim.Disappear(gtx.Now)
						}
					case component.AppBarContextMenuDismissed:
						log.Printf("Context menu dismissed: %v", event)
					case component.AppBarOverflowActionClicked:
						log.Printf("Overflow action selected: %v", event)
					}
				}
				if alternatePalette.Value {
					th.Palette = altPalette
					currentAccent = altPaletteAccent
				} else {
					th.Palette = lightPalette
					currentAccent = lightPaletteAccent
				}
				if nav.NavDestinationChanged() {
					page := pages[nav.CurrentNavDestination().(int)]
					bar.Title = page.Name
					bar.SetActions(page.Actions, page.Overflow)
				}
				paint.Fill(gtx.Ops, th.Palette.Bg)
				layout.Inset{
					Top:    e.Insets.Top,
					Bottom: e.Insets.Bottom,
					Left:   e.Insets.Left,
					Right:  e.Insets.Right,
				}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					content := layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								gtx.Constraints.Max.X /= 3
								return nav.Layout(gtx, th, &navAnim)
							}),
							layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
								page := &pages[nav.CurrentNavDestination().(int)]
								return page.List.Layout(gtx, 1, func(gtx C, _ int) D {
									return layout.UniformInset(unit.Dp(4)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
										return page.layout(gtx)
									})
								})
							}),
						)
					})
					bar := layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return bar.Layout(gtx, th)
					})
					flex := layout.Flex{Axis: layout.Vertical}
					if bottomBar.Value {
						flex.Layout(gtx, content, bar)
					} else {
						flex.Layout(gtx, bar, content)
					}
					modal.Layout(gtx, th)
					return layout.Dimensions{Size: gtx.Constraints.Max}
				})
				e.Frame(gtx.Ops)
			}
		}
	}
}
