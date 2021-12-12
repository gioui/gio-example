package ui

import (
	"gioui.org/app"
	"gioui.org/example/video/packet"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"image"
	"time"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

type CurrentView int

const (
	StateView CurrentView = iota
	VideoView
)

const DefaultAudioOffset = -1
const VisibilityTimeOut = int64(3000)

type Player struct {
	*app.Window
	ticker                     *time.Ticker
	th                         *material.Theme
	perSecond                  <-chan time.Time
	uiSeekerWidget             widget.Float
	currentImage               image.Image
	currentView                CurrentView
	player                     *packet.Controller
	filePath                   string
	status                     PlayerStatus
	audios                     packet.Audios
	audioBuffChan              chan [2]float64
	stateList                  layout.List
	timeChan                   chan int64
	video                      *packet.Video
	videoChan                  chan *packet.Video
	playBtn                    widget.Clickable
	pauseBtn                   widget.Clickable
	stopBtn                    widget.Clickable
	showStateBtn               widget.Clickable
	openFileBtn                widget.Clickable
	screenClickable            widget.Clickable
	rewindButton               widget.Clickable
	forwardButton              widget.Clickable
	incrementButton            widget.Clickable
	decrementButton            widget.Clickable
	initialized                bool
	isPathValid                bool
	speakerInitialized         bool
	isFullScreen               bool
	RightKeyPressed            bool
	UpKeyPressed               bool
	LeftKeyPressed             bool
	DownKeyPressed             bool
	SpaceKeyPressed            bool
	EscKeyPressed              bool
	sampleRate                 int
	frameIndex                 int
	audioOffset                int64
	lastHoveredTime            int64
	lastClickedTime            int64
	RightKeyLastUpdated        int64
	UpKeyLastUpdated           int64
	LeftKeyLastUpdated         int64
	DownKeyLastUpdated         int64
	SpaceKeyLastUpdated        int64
	EscKeyLastUpdated          int64
	forwardButtonLastUpdated   int64
	rewindButtonLastUpdated    int64
	incrementButtonLastUpdated int64
	decrementButtonLastUpdated int64
	sliderLastUpdated          int64
}

func (p *Player) CurrentView() CurrentView {
	return p.currentView
}
func (p *Player) SetCurrentView(v CurrentView) {
	p.currentView = v
}

func (p *Player) OnVideoPathChanged() {
	if p.IsPathValid() {
		go p.player.LoadFile(p.FilePath())
		p.audioOffset = DefaultAudioOffset
	}
}

func (p *Player) Layout(gtx C) D {
	if p.th == nil {
		p.th = material.NewTheme(gofont.Collection())
		p.th.Bg, p.th.Fg, p.th.ContrastBg, p.th.ContrastFg = p.th.Fg, p.th.Bg, p.th.Fg, p.th.Bg
	}

	if !p.initialized {
		p.initialize()
	}
	p.update(gtx)
	p.handleKeyboardArrowsEvent()
	p.handleUIEvents(gtx)
	p.drawHover(gtx)

	d := layout.Stack{
		Alignment: layout.S,
	}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
			fl := layout.Flex{Axis: layout.Vertical}
			d := fl.Layout(gtx,
				layout.Flexed(1, func(gtx C) D {
					switch p.CurrentView() {
					default:
						fallthrough
					case StateView:
						return p.drawStateView(gtx, p.th)
					case VideoView:
						return p.drawVideoView(gtx, p.currentImage)
					}
				}),
			)
			p.drawVideoOverlay(gtx)
			return d
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			if time.Now().UnixMilli()-p.lastHoveredTime < VisibilityTimeOut ||
				p.CurrentView() == StateView {
				return p.drawFooter(gtx)
			}
			return D{Size: gtx.Constraints.Max}
		}),
	)
	op.InvalidateOp{}.Add(gtx.Ops)
	return d
}
