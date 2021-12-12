package ui

import (
	"gioui.org/example/video/packet"
	"gioui.org/layout"
	"time"
)

func (p *Player) initialize() {
	if p.perSecond == nil {
		p.perSecond = time.Tick(time.Second)
	}
	if p.player == nil {
		var filePath string
		if ok := p.IsPathValid(); ok {
			filePath = p.FilePath()
		}
		if p.timeChan == nil {
			p.timeChan = make(chan int64, 1)
		}
		if p.videoChan == nil {
			p.videoChan = make(chan *packet.Video, 1)
		}

		p.player = packet.New(filePath, p.timeChan, p.videoChan, func(audios packet.Audios) {
			p.audios = audios
			p.SetCurrentView(VideoView)
			p.status = Playing
		})
	}
	if p.ticker == nil {
		var framesCount int
		if p.video != nil {
			framesCount = len(p.video.Images())
		}
		if framesCount == 0 {
			framesCount = 30
		}
		p.ticker = time.NewTicker(time.Duration(float64(time.Second) / float64(framesCount)))
	}
	p.stateList.Axis = layout.Vertical
	p.stateList.Alignment = layout.Middle
	p.initialized = true
}
