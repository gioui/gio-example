package packet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/zergon321/reisen"
	"image"
	"log"
	"time"
)

func (p *Controller) decodeVideo(done <-chan bool) {
	var status Status
	var message string

	defer func() {
		if r := recover(); r != nil {
			status = Error
			message = fmt.Sprintf("%s", r)
			p.decodingVideoState.setState(status, message)
		}
		if status == Running {
			status = Error
			message = "Unknown Error"
		}
		p.decodingVideoState.setState(status, message)
	}()

	status = Running
	message = "Decoding Video Running..."
	p.decodingVideoState.setState(status, message)
	var t int64
	var videos = make(Videos)
	var ignoreVideos = make(Videos)

	for {
		select {
		case <-done:
			status = Canceled
			message = "Cancelled"
			p.decodingVideoState.setState(status, message)
			return
		case nt, ok := <-p.timeChan:
			if ok {
				t = nt
			}
			startPos := (t / VideoSize) * VideoSize
			if startPos > p.EndTime() || startPos < 0 {
				message = fmt.Sprintf("%d startPos out of bounds", startPos)
				p.decodingVideoState.setState(status, message)
				continue
			}

			vid, ok := videos[t]
			if ok {
				select {
				case p.videoChan <- vid:
				default:
				}
				_, nxtOk := videos[t+VideoSizeOffset]
				if nxtOk {
					continue
				}
				newVideos := make(Videos)
				ignoreVideos = make(Videos)
				for i := t; i < t+VideoSizeOffset; i++ {
					if newVid, ok := videos[i]; ok {
						newVideos[i] = newVid
						ignoreVideos[i] = newVid
					}
				}
				videos = newVideos
			}

			if !ok {
				videos = make(Videos)
			}

			_, ok = videos[t]
			// check if it's the last package and video exist at current position
			endPos := (p.EndTime() / VideoSize) * VideoSize
			if startPos == endPos && ok {
				continue
			}

			if startPos > p.EndTime() {
				continue
			}

			videoPath, ok := p.videosPath[startPos]
			if !ok {
				message = fmt.Sprintf("path %s doesn't exist, this indicates bug in the app...",
					videoPath,
				)
				status = Error
				return
			}
			media, err := reisen.NewMedia(videoPath)
			if err != nil {
				status = Error
				message = err.Error()
				return
			}
			err = media.OpenDecode()
			if err != nil {
				status = Error
				message = err.Error()
				return
			}
			videoStream := media.VideoStreams()[0]
			err = videoStream.Open()
			if err != nil {
				status = Error
				message = err.Error()
				return
			}
			gotVideo := true
			var pkt *reisen.Packet

			// Decoding process starts here...
			message = "current videos decoding started..."
			p.decodingVideoState.setState(status, message)
			for gotVideo {
				select {
				case <-done:
					status = Canceled
					message = "Cancelled"
					videos = nil
					_ = videoStream.Close()
					_ = media.CloseDecode()
					media.Close()
					return
				case nt, ok = <-p.timeChan:
					if ok {
						if t != nt {
							t = nt
							if vid, ok = videos[t]; ok {
								select {
								case p.videoChan <- vid:
								default:
								}
							}
						}
					}
					break
				default:
				}
				pkt, gotVideo, err = media.ReadPacket()
				if !gotVideo {
					break
				}
				if err != nil {
					p.decodingVideoState.setState(status, err.Error())
				}

				nextPos := startPos + VideoSize + VideoSizeOffset
				if t < startPos || t > nextPos {
					break
				}

				switch pkt.Type() {
				case reisen.StreamVideo:
					s := media.Streams()[pkt.StreamIndex()].(*reisen.VideoStream)
					videoFrame, gotFrame, err := s.ReadVideoFrame()
					if !gotFrame {
						break
					}
					if videoFrame == nil {
						continue
					}

					if err != nil {
						status = Error
						message = err.Error()
						return
					}

					d, err := videoFrame.PresentationOffset()

					if err != nil {
						status = Error
						message = err.Error()
						return
					}

					if err == nil {
						b := videoFrame.Image().Bounds()
						img := videoFrame.Image().SubImage(b)
						ti := int64(d.Seconds()) + startPos

						if _, ok := ignoreVideos[ti]; ok {
							continue
						}
						if video, ok := videos[ti]; ok {
							video.appendImage(img)
						} else {
							videos[ti] = &Video{time: ti, images: []image.Image{img}}
						}
					} else {
						p.decodingVideoState.setState(status, err.Error())
					}
				}
			}
			_ = videoStream.Close()
			_ = media.CloseDecode()
			media.Close()
			err = p.checkVideoSize(videos)
			if err != nil {
				status = Error
				message = err.Error()
				return
			}
			message = "current videos decoding ended..."
			p.decodingVideoState.setState(status, message)
		default:

		}
	}
}

func (p *Controller) decodeAudio(done <-chan bool) {
	var status Status
	var message string
	start := time.Now().Unix()
	defer func() {
		if r := recover(); r != nil {
			status = Error
			message = fmt.Sprintf("%s", r)
			p.decodingAudioState.setState(status, message)
		}
		if status == Running {
			status = Error
			message = "Unknown Error"
		}
		p.decodingAudioState.setState(status, message)
	}()

	status = Running
	message = "Decoding..."
	p.decodingAudioState.setState(status, message)

	select {
	case <-done:
		status = Canceled
		message = "Cancelled"
		return
	default:
		audioPath := p.audioPath
		media, err := reisen.NewMedia(audioPath)
		if err != nil {
			status = Error
			message = err.Error()
			return
		}
		err = media.OpenDecode()
		if err != nil {
			status = Error
			message = err.Error()
			return
		}
		audioStream := media.Streams()[0].(*reisen.AudioStream)
		err = audioStream.Open()
		if err != nil {
			status = Error
			message = err.Error()
			return
		}
		gotAudio := true
		var pkt *reisen.Packet

		p.audios = make(Audios)

		for gotAudio {
			select {
			case <-done:
				status = Canceled
				message = "Cancelled"
				return
			default:
				pkt, gotAudio, err = media.ReadPacket()
				if !gotAudio {
					break
				}
				if err != nil {
					p.decodingAudioState.setState(status, err.Error())
				}

				switch pkt.Type() {
				case reisen.StreamAudio:
					s := media.Streams()[pkt.StreamIndex()].(*reisen.AudioStream)
					audioFrame, gotFrame, _ := s.ReadAudioFrame()

					if !gotFrame {
						gotAudio = false
						break
					}

					if audioFrame == nil {
						continue
					}

					d, err := audioFrame.PresentationOffset()

					if err == nil {
						pack := p.audios[int64(d.Seconds())]
						reader := bytes.NewReader(audioFrame.Data())
						for reader.Len() > 0 {
							sample := [2]float64{0, 0}
							var result float64
							err := binary.Read(reader, binary.LittleEndian, &result)

							if err != nil {
								log.Println(err)
							}
							sample[0] = result

							err = binary.Read(reader, binary.LittleEndian, &result)

							if err != nil {
								log.Println(err)
							}

							sample[1] = result
							pack = append(pack, sample)
						}
						p.audios[(int64(d.Seconds()))] = pack
					} else {
						p.decodingAudioState.setState(status, err.Error())
					}
				}
			}
		}
		_ = audioStream.Close()
		_ = media.CloseDecode()
		media.Close()
		status = Success
		end := time.Now().Unix() - start
		message = fmt.Sprintf("Time taken: %s", (time.Duration(end) * time.Second).String())
		return
	}

}

func (p *Controller) checkVideoSize(videos Videos) error {
	if (int64(len(videos))) > (VideoSize + (VideoSizeOffset * 2)) {
		err := errors.New(fmt.Sprintf("videos size exceeded %d, this indicate serious bug in code, that may cause"+
			" freeze OS by consuming high memory...", len(videos)))
		return err
	}
	return nil
}
