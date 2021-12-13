package ui

import (
	"gioui.org/app"
	"gioui.org/io/pointer"
	"gioui.org/op"
	"github.com/sqweek/dialog"
	"math"
	"time"
)

const RightKeyInterval = 10
const LeftKeyInterval = 10
const UpKeyInterval = 5
const DownKeyInterval = 5

type KeyAction int

const (
	ActionIncrement KeyAction = iota
	ActionDecrement
)

func (p *Player) keyboardSeekHelper(lastUpdatedAt *int64, interval int64, keyActive bool, keyAction KeyAction) {
	now := time.Now().UnixMilli()
	diff := now - *lastUpdatedAt
	if keyActive {
		p.lastHoveredTime = time.Now().UnixMilli()
	}
	if keyActive && *lastUpdatedAt != 0 {
		if diff >= interval {
			val := p.uiSeekerWidget.Value
			if keyAction == ActionIncrement {
				val += 1
			}
			if keyAction == ActionDecrement {
				val -= 1
			}
			val = float32(math.Min(float64(val), float64(p.player.EndTime())))
			val = float32(math.Max(float64(val), 0))
			p.uiSeekerWidget.Value = val
			*lastUpdatedAt = 0
		}
	}
}

func (p *Player) handleKeyboardArrowsEvent() {
	p.keyboardSeekHelper(&p.UpKeyLastUpdated, UpKeyInterval, p.UpKeyPressed, ActionIncrement)
	p.keyboardSeekHelper(&p.RightKeyLastUpdated, RightKeyInterval, p.RightKeyPressed, ActionIncrement)
	p.keyboardSeekHelper(&p.DownKeyLastUpdated, DownKeyInterval, p.DownKeyPressed, ActionDecrement)
	p.keyboardSeekHelper(&p.LeftKeyLastUpdated, LeftKeyInterval, p.LeftKeyPressed, ActionDecrement)
}

func (p *Player) handleUIEvents(gtx C) {
	if p.uiSeekerWidget.Dragging() {
		p.sliderLastUpdated = time.Now().UnixMilli()
		p.lastHoveredTime = time.Now().UnixMilli()
	}

	if p.playBtn.Clicked() {
		if !p.IsPathValid() {
			p.SetCurrentView(StateView)
		}
		if p.status != Playing {
			p.status = Playing
		}
	}

	if p.pauseBtn.Clicked() {
		if p.status != Paused {
			p.status = Paused
		}
	}

	if p.showStateBtn.Clicked() {
		switch p.CurrentView() {
		default:
			fallthrough
		case StateView:
			p.SetCurrentView(VideoView)
		case VideoView:
			p.SetCurrentView(StateView)
		}
	}

	if p.stopBtn.Clicked() {
		p.uiSeekerWidget.Value = 0
		p.video = nil
		p.currentImage = nil
		p.status = Stopped
	}

	if p.SpaceKeyPressed && p.SpaceKeyLastUpdated != 0 {
		p.lastHoveredTime = time.Now().UnixMilli()
		p.SpaceKeyLastUpdated = 0
		if p.status == Playing {
			p.status = Paused
		} else if p.status == Paused {
			p.status = Playing
		}
	}

	if p.EscKeyPressed && p.EscKeyLastUpdated != 0 {
		p.lastHoveredTime = time.Now().UnixMilli()
		p.EscKeyLastUpdated = 0
		if p.isFullScreen {
			p.Window.Option(app.Windowed.Option())
		}
	}

	for _, ev := range gtx.Events(&p.lastHoveredTime) {
		switch ev := ev.(type) {
		case pointer.Event:
			switch ev.Type {
			case pointer.Move:
				p.lastHoveredTime = time.Now().UnixMilli()
			}
		}
	}

	if p.openFileBtn.Clicked() {
		filePath, err := dialog.File().Filter("Video file").Load()
		if err != nil {
			if err.Error() != "Cancelled" {
				p.isPathValid = false
				p.filePath = filePath
				p.OnVideoPathChanged()
			}
		} else {
			if p.filePath != filePath {
				p.filePath = filePath
				p.isPathValid = true
				p.OnVideoPathChanged()
			}
		}
	}

	if p.screenClickable.Clicked() {
		// check if doubleClicked
		if time.Now().UnixMilli()-p.lastClickedTime < 350 {
			if p.isFullScreen {
				if p.status == Playing {
					p.status = Paused
					go func() {
						p.Window.Option(app.Windowed.Option())
						p.status = Playing
					}()
				} else {
					p.Window.Option(app.Windowed.Option())
				}
			} else {
				p.Window.Option(app.Fullscreen.Option())
			}
			p.isFullScreen = !p.isFullScreen
			op.InvalidateOp{}.Add(gtx.Ops)
		}
		p.lastClickedTime = time.Now().UnixMilli()
	}

	if p.incrementButton.Clicked() {
		now := time.Now().UnixMilli()
		diff := now - p.incrementButtonLastUpdated
		if diff < 350 {
			p.audioOffset++
			p.incrementButtonLastUpdated = 0
		}
	}
	if p.decrementButton.Clicked() {
		now := time.Now().UnixMilli()
		diff := now - p.decrementButtonLastUpdated
		if diff < 350 {
			p.audioOffset--
			p.decrementButtonLastUpdated = 0
		}
	}

	if p.forwardButton.Pressed() {
		if p.forwardButtonLastUpdated == 0 {
			p.forwardButtonLastUpdated = time.Now().UnixMilli()
		}
	}

	if p.rewindButton.Pressed() {
		if p.rewindButtonLastUpdated == 0 {
			p.rewindButtonLastUpdated = time.Now().UnixMilli()
		}
	}

	p.keyboardSeekHelper(&p.forwardButtonLastUpdated, RightKeyInterval, p.forwardButton.Pressed(), ActionIncrement)
	p.keyboardSeekHelper(&p.rewindButtonLastUpdated, LeftKeyInterval, p.rewindButton.Pressed(), ActionDecrement)

	if p.incrementButton.Pressed() {
		if p.incrementButtonLastUpdated == 0 {
			p.incrementButtonLastUpdated = time.Now().UnixMilli()
		}
		now := time.Now().UnixMilli()
		diff := now - p.incrementButtonLastUpdated
		if diff > 350 {
			p.audioOffset += RightKeyInterval
			p.incrementButtonLastUpdated = 0
		}
		p.lastHoveredTime = time.Now().UnixMilli()
	}

	if p.decrementButton.Pressed() {
		if p.decrementButtonLastUpdated == 0 {
			p.decrementButtonLastUpdated = time.Now().UnixMilli()
		}
		now := time.Now().UnixMilli()
		diff := now - p.decrementButtonLastUpdated
		if diff > 350 {
			p.audioOffset -= LeftKeyInterval
			p.decrementButtonLastUpdated = 0
		}
		p.lastHoveredTime = time.Now().UnixMilli()
	}
}
