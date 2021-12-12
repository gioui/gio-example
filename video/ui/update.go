package ui

import (
	"image"
)

type PlayerStatus int

const (
	Idle PlayerStatus = iota
	Loading
	Playing
	Paused
	Stopped
)

func (p *Player) canPlay() bool {
	return p.status == Playing &&
		!p.rewindButton.Pressed() &&
		!p.forwardButton.Pressed() &&
		!p.RightKeyPressed &&
		!p.LeftKeyPressed &&
		!p.UpKeyPressed &&
		!p.DownKeyPressed &&
		!p.uiSeekerWidget.Dragging()
}

func (p *Player) update(gtx C) {
	if p.canPlay() {
		p.response()
		p.request()
		select {
		case <-p.ticker.C:
			t := int64(p.uiSeekerWidget.Value)
			if p.video != nil && p.video.Time() == t {
				frameSize := len(p.video.Images())
				if p.frameIndex < frameSize {
					p.currentImage = p.video.Images()[p.frameIndex]
					if p.speakerInitialized {
						buff, ok := p.audios[t-p.audioOffset]
						if ok {
							buffSize := len(buff)
							buffPerFrame := buffSize / frameSize
							startIndex := buffPerFrame * p.frameIndex
							endIndex := buffPerFrame * (p.frameIndex + 1)
							reqBuff := buff[startIndex:endIndex]
							for _, eachBuff := range reqBuff {
								p.audioBuffChan <- eachBuff
							}
						}
					}
					p.frameIndex++
				}
			}
		default:
		}
	}

	if p.currentImage == nil {
		rgb := &image.RGBA{Pix: nil, Stride: 0, Rect: image.Rectangle{}}
		p.currentImage = rgb.SubImage(rgb.Bounds())
	}
}
