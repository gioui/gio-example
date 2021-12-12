package packet

import (
	"sync"
)

type Status int

const (
	Idle Status = iota
	Error
	Running
	Canceled
	Success
)

func (s Status) String() string {
	switch s {
	case Idle:
		return "Idle"
	case Error:
		return "Error"
	case Running:
		return "Running"
	case Canceled:
		return "Canceled"
	case Success:
		return "Success"
	}
	return "Unknown"
}

type State interface {
	Status() Status
	Message() string
}

type Audio [][2]float64
type Audios map[int64]Audio

// VideoSize represents videos size in seconds as well as number of videos
// (The video files are split and named based on this value)
// It's a key value for the app-
const VideoSize = int64(10)

const VideoSizeOffset = int64(2)

type CommonState struct {
	stateMutex sync.RWMutex
	status     Status
	message    string
}

func (d *CommonState) Status() Status {
	d.stateMutex.RLock()
	defer d.stateMutex.RUnlock()
	return d.status
}

func (d *CommonState) Message() string {
	d.stateMutex.RLock()
	defer d.stateMutex.RUnlock()
	return d.message
}

func (d *CommonState) setState(status Status, message string) {
	d.stateMutex.Lock()
	defer d.stateMutex.Unlock()
	d.status = status
	d.message = message
}
