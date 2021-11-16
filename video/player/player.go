package video

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/zergon321/reisen"
	"image"
	"log"
	"time"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

const (
	frameBufferSize                   = 1024
	sampleRate                        = 44100
	channelCount                      = 2
	bitDepth                          = 8
	sampleBufferSize                  = 32 * channelCount * bitDepth * 1024
	SpeakerSampleRate beep.SampleRate = 44100
)

type Player struct {
	image                  *image.Image
	audioSource            <-chan [2]float64
	ticker                 <-chan time.Time
	errs                   <-chan error
	frameBuffer            <-chan *image.RGBA
	fps                    int
	videoTotalFramesPlayed int
	videoPlaybackFPS       int
	perSecond              <-chan time.Time
	last                   time.Time
	deltaTime              float64
	Filepath               string
	th                     *material.Theme
}

func readVideoAndAudio(media *reisen.Media) (<-chan *image.RGBA, <-chan [2]float64, chan error, error) {
	frameBuffer := make(chan *image.RGBA,
		frameBufferSize)
	sampleBuffer := make(chan [2]float64, sampleBufferSize)
	errs := make(chan error)

	err := media.OpenDecode()

	if err != nil {
		return nil, nil, nil, err
	}

	videoStream := media.VideoStreams()[0]
	err = videoStream.Open()

	if err != nil {
		return nil, nil, nil, err
	}

	audioStream := media.AudioStreams()[0]
	err = audioStream.Open()

	if err != nil {
		return nil, nil, nil, err
	}

	go func() {
		for {
			packet, gotPacket, err := media.ReadPacket()

			if err != nil {
				go func(err error) {
					errs <- err
				}(err)
			}

			if !gotPacket {
				break
			}

			switch packet.Type() {
			case reisen.StreamVideo:
				s := media.Streams()[packet.StreamIndex()].(*reisen.VideoStream)
				videoFrame, gotFrame, err := s.ReadVideoFrame()

				if err != nil {
					go func(err error) {
						errs <- err
					}(err)
				}

				if !gotFrame {
					break
				}

				if videoFrame == nil {
					continue
				}

				frameBuffer <- videoFrame.Image()

			case reisen.StreamAudio:
				s := media.Streams()[packet.StreamIndex()].(*reisen.AudioStream)
				audioFrame, gotFrame, err := s.ReadAudioFrame()

				if err != nil {
					go func(err error) {
						errs <- err
					}(err)
				}

				if !gotFrame {
					break
				}

				if audioFrame == nil {
					continue
				}

				// Turn the raw byte data into
				// audio samples of type [2]float64.
				reader := bytes.NewReader(audioFrame.Data())

				for reader.Len() > 0 {
					sample := [2]float64{0, 0}
					var result float64
					err = binary.Read(reader, binary.LittleEndian, &result)

					if err != nil {
						go func(err error) {
							errs <- err
						}(err)
					}

					sample[0] = result

					err = binary.Read(reader, binary.LittleEndian, &result)

					if err != nil {
						go func(err error) {
							errs <- err
						}(err)
					}

					sample[1] = result
					sampleBuffer <- sample
				}
			}
		}

		videoStream.Close()
		audioStream.Close()
		media.CloseDecode()
		close(frameBuffer)
		close(sampleBuffer)
		close(errs)
	}()

	return frameBuffer, sampleBuffer, errs, nil
}

func audioStreamers(sampleSource <-chan [2]float64) beep.Streamer {
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

func (p *Player) Start() (err error) {
	// Initialize the audio speaker.
	err = speaker.Init(sampleRate, SpeakerSampleRate.N(time.Second/10))

	if err != nil {
		log.Println(err)
		return
	}

	media, err := reisen.NewMedia(p.Filepath)
	if err != nil {
		log.Println(err)
		return
	}

	// Get the FPS for playing
	// video frames.
	videoFPS, _ := media.Streams()[0].FrameRate()

	if err != nil {
		log.Println(err)
		return
	}

	// SPF for frame ticker.
	spf := 1.0 / float64(videoFPS)
	frameDuration, err := time.
		ParseDuration(fmt.Sprintf("%fs", spf))

	if err != nil {
		log.Println(err)
		return
	}

	p.frameBuffer, p.audioSource,
		p.errs, err = readVideoAndAudio(media)

	if err != nil {
		return
	}

	// Start playing audio samples.
	speaker.Play(audioStreamers(p.audioSource))

	p.ticker = time.Tick(frameDuration)

	// Setup metrics.
	p.last = time.Now()
	p.fps = 0
	p.perSecond = time.Tick(time.Second)
	p.videoTotalFramesPlayed = 0
	p.videoPlaybackFPS = 0
	return err
}

func (p *Player) update(gtx C) {
	// Compute dt.
	p.deltaTime = time.Since(p.last).Seconds()
	p.last = time.Now()

	// Check for incoming errors.
	select {
	case err, ok := <-p.errs:
		if ok {
			log.Println(err)
		}

	default:
	}

	select {
	case <-p.ticker:
		frame, ok := <-p.frameBuffer

		if ok {
			b := frame.Bounds()
			img := frame.SubImage(b)
			p.image = &img
			p.videoTotalFramesPlayed++
			p.videoPlaybackFPS++
		}

	default:
	}

	p.fps++

	select {
	case <-p.perSecond:
		log.Printf("%s | FPS: %d | dt: %f | Frames: %d | Video FPS: %d",
			"Video", p.fps, p.deltaTime, p.videoTotalFramesPlayed, p.videoPlaybackFPS)

		p.fps = 0
		p.videoPlaybackFPS = 0
	default:
	}
	op.InvalidateOp{}.Add(gtx.Ops)
}

func (p *Player) Layout(gtx C) D {
	if p.th == nil {
		p.th = material.NewTheme(gofont.Collection())
	}
	p.update(gtx)
	return layout.Inset{
		Top:    unit.Dp(32),
		Right:  unit.Dp(32),
		Bottom: unit.Dp(32),
		Left:   unit.Dp(32),
	}.Layout(gtx,
		func(gtx C) D {
			fl := layout.Flex{
				Axis:      layout.Vertical,
				Spacing:   0,
				Alignment: 0,
				WeightSum: 0,
			}
			return fl.Layout(gtx,
				layout.Flexed(1, func(gtx C) D {
					if p.image == nil {
						return layout.Center.Layout(gtx, material.Body1(p.th, "No Video Is Currently Playing").Layout)
					}
					imgOps := paint.NewImageOp(*p.image)
					return widget.Image{
						Src:      imgOps,
						Fit:      widget.Fill,
						Position: layout.Center,
						Scale:    0,
					}.Layout(gtx)
				}),
			)
		},
	)
}
