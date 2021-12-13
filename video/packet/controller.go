package packet

import (
	"log"
	"sync"
)

type Controller struct {
	audios                 Audios
	audioBuffChan          chan [2]float64
	filePath               string
	encodingState          CommonState
	decodingVideoState     CommonState
	decodingAudioState     CommonState
	startTime              int64
	endTime                int64
	endTimeMutex           sync.RWMutex
	videosPath             map[int64]string
	audioPath              string
	done                   chan bool
	timeChan               <-chan int64
	videoChan              chan<- *Video
	OnAudioDecodingSuccess func(audios Audios)
}

func New(filePath string,
	timeChan <-chan int64,
	videoChan chan<- *Video,
	OnAudioDecodingSuccess func(audios Audios)) *Controller {
	p := &Controller{
		filePath:               filePath,
		timeChan:               timeChan,
		videoChan:              videoChan,
		OnAudioDecodingSuccess: OnAudioDecodingSuccess,
	}
	if timeChan == nil || videoChan == nil {
		log.Fatalln("timeChan and videoChan are required and cannot be nil")
	}
	if filePath != "" {
		go p.LoadFile(filePath)
	}
	return p
}

func (p *Controller) LoadFile(filePath string) {
	if filePath == "" {
		return
	}
	if p.timeChan == nil || p.videoChan == nil {
		log.Fatalln("timeChan and videoChan are required and cannot be nil")
	}
	go func() {
		if p.done != nil {
			select {
			case p.done <- true:
			default:
			}
			close(p.done)
		}
		p.done = make(chan bool, 10)
		p.filePath = filePath
		p.videosPath = make(map[int64]string, 0)
		p.encodingState.setState(Idle, "Idle")
		p.decodingVideoState.setState(Idle, "Idle")
		p.decodingAudioState.setState(Idle, "Idle")
		p.encode()
		if p.encodingState.Status() != Success {
			return
		}
		p.decodeAudio(p.done)
		if p.decodingAudioState.Status() != Success {
			return
		}
		if p.OnAudioDecodingSuccess != nil {
			p.OnAudioDecodingSuccess(p.audios)
		}
		go p.decodeVideo(p.done)
	}()
}
