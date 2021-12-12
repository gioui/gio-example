package packet

import (
	"image"
	"sync"
)

type Video struct {
	videoMutex sync.RWMutex
	time       int64
	images     []image.Image
}

type Videos map[int64]*Video

func (v *Video) Time() int64 {
	v.videoMutex.RLock()
	defer v.videoMutex.RUnlock()
	return v.time
}

func (v *Video) Images() []image.Image {
	v.videoMutex.RLock()
	defer v.videoMutex.RUnlock()
	return v.images
}

func (v *Video) clearImages() {
	v.videoMutex.Lock()
	defer v.videoMutex.Unlock()
	v.images = make([]image.Image, 0)
}

func (v *Video) appendImages(images []image.Image) {
	v.videoMutex.Lock()
	defer v.videoMutex.Unlock()
	v.images = append(v.images, images...)
}

func (v *Video) appendImage(image image.Image) {
	v.videoMutex.Lock()
	defer v.videoMutex.Unlock()
	v.images = append(v.images, image)
}
