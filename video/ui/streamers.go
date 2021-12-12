package ui

import (
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"log"
	"time"
)

func AudioStreamers(sampleSource <-chan [2]float64) beep.Streamer {
	return beep.StreamerFunc(func(samples [][2]float64) (n int, ok bool) {
		numRead := 0

		for i := 0; i < len(samples); i++ {
			sample, ok := <-sampleSource

			if !ok {
				numRead = i + 1
				break
			}

			samples[i] = sample
			numRead++
		}

		if numRead < len(samples) {
			return numRead, false
		}

		return numRead, true
	})
}

func (p *Player) InitializeSpeakers() {
	speaker.Clear()
	speaker.Close()
	p.speakerInitialized = false

	if p.video == nil {
		return
	}
	t := int64(p.uiSeekerWidget.Value)
	audios, _ := p.audios[t]
	p.sampleRate = len(audios)
	if p.sampleRate == 0 {
		return
	}
	err := speaker.Init(beep.SampleRate(p.sampleRate), beep.SampleRate(p.sampleRate).N(time.Second))
	if err != nil {
		log.Println(err)
	} else {
		log.Println("speaker initialized....")
		p.speakerInitialized = true
	}
	channelCount := 2
	bitDepth := 8
	sampleBufferSize := 32 * channelCount * bitDepth * 1024 * 2
	if p.audioBuffChan == nil {
		p.audioBuffChan = make(chan [2]float64, sampleBufferSize)
	}
	speaker.Play(AudioStreamers(p.audioBuffChan))
}
