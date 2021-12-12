package ui

import (
	"time"
)

func (p *Player) response() {
	select {
	case video := <-p.videoChan:
		t := int64(p.uiSeekerWidget.Value)
		if video == nil || video.Time() != t {
			break
		}
		p.video = video
		p.frameIndex = 0
		framesCount := len(video.Images())
		if framesCount != 0 {
			p.ticker.Reset(time.Second / time.Duration(framesCount))
		}
		if !p.speakerInitialized {
			p.InitializeSpeakers()
		}

	default:

	}
}
