package packet

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const AppDirName = "gioplayer"
const VideosDirName = "videos"
const CachedFileName = "filepath"
const AudioFileName = "audio.aac"

func (p *Controller) FilePath() string {
	return p.filePath
}

func (p *Controller) VideoDecodingState() *CommonState {
	return &p.decodingVideoState
}

func (p *Controller) AudioDecodingState() *CommonState {
	return &p.decodingAudioState
}

func (p *Controller) EncodingState() *CommonState {
	return &p.encodingState
}

func (p *Controller) StartTime() int64 {
	return p.startTime
}

func (p *Controller) EndTime() int64 {
	p.endTimeMutex.RLock()
	defer p.endTimeMutex.RUnlock()
	return p.endTime
}

func (p *Controller) SetStartTime(s int64) {
	p.startTime = s
}

func (p *Controller) SetEndTime(s int64) {
	p.endTimeMutex.Lock()
	defer p.endTimeMutex.Unlock()
	p.endTime = s
}

func (p *Controller) encode() {
	var status Status
	var message string
	defer func() {
		if r := recover(); r != nil {
			status = Error
			message = fmt.Sprintf("%s", r)
			p.encodingState.setState(status, message)
		}
		if status == Running {
			status = Error
			message = "Unknown Error"
		}
		p.encodingState.setState(status, message)
	}()

	status = Running
	message = "Encoding started..."
	p.encodingState.setState(status, message)

	dir := os.TempDir()
	dir = filepath.Join(dir, AppDirName, VideosDirName)

	cachedPath := filepath.Join(dir, CachedFileName)
	audioPath := filepath.Join(dir, AudioFileName)
	p.audioPath = audioPath

	_, err := os.Stat(cachedPath)
	fileExist := !errors.Is(err, os.ErrNotExist)

	if fileExist {
		file, err := os.Open(cachedPath)
		if err == nil {
			b, err := ioutil.ReadAll(file)
			if err == nil {
				previousFilePath := string(b)
				if previousFilePath != p.FilePath() {
					fileExist = false
				}
			}
		}
		err = file.Close()
		if err != nil {
			fileExist = false
		} else {
			message = fmt.Sprintf("File exist, using cached path %s", cachedPath)
			p.encodingState.setState(status, message)
		}
	}

	if !fileExist {
		err = os.RemoveAll(dir)
		if err != nil {
			message = err.Error()
			status = Error
			return
		}
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			message = err.Error()
			status = Error
			return
		}
	}

	argSlice := []string{"-i", p.FilePath(), "-show_entries", "format=duration", "-v", "quiet", "-of", "csv=p=0"}
	command := exec.Command("ffprobe", argSlice...)
	b, err := command.CombinedOutput()
	if err != nil {
		message = err.Error()
		status = Error
		return
	}

	endTimeStr := strings.Split(string(b), "\n")[0]
	if runtime.GOOS == "windows" {
		endTimeStr = strings.Split(string(b), "\r\n")[0]
	}
	endTimeFl, err := strconv.ParseFloat(endTimeStr, 64)
	if err != nil {
		message = err.Error()
		status = Error
		return
	}
	endTime := int64(endTimeFl)

	p.SetEndTime(endTime)

	command = nil
	start := time.Now().Unix()
	for i := int64(0); i <= endTime; i += VideoSize {
		videoPath := filepath.Join(dir, fmt.Sprintf("%d.mp4", i))
		currEndTime := VideoSize
		if endTime-i < VideoSize {
			currEndTime = endTime - i
		}

		if !fileExist {
			// alternate option if -to flag is used instead of t
			//currEndTime := int(math.Min(float64(p.EndTime()), float64(i + interval)))
			argSlice := []string{
				"-ss", fmt.Sprintf("%d", i),
				"-i", p.FilePath(),
				"-t", fmt.Sprintf("%d", currEndTime+VideoSizeOffset),
				"-c:v", "copy",
				videoPath,
			}
			cmd := exec.Command("ffmpeg", argSlice...)
			_, err := cmd.CombinedOutput()
			if err != nil {
				message = err.Error()
				status = Error
				return
			}
			p.encodingState.setState(status, fmt.Sprintf("ffmpeg %s", strings.Join(argSlice, " ")))
		}
		p.videosPath[i] = videoPath
	}
	end := time.Now().Unix()
	dur := time.Duration(end-start) * time.Second
	if !fileExist {
		// alternate option if -to flag is used instead of t
		//currEndTime := int(math.Min(float64(p.EndTime()), float64(i + interval)))
		argSlice := []string{
			"-ss", fmt.Sprintf("%d", 0),
			"-i", p.FilePath(),
			"-t", fmt.Sprintf("%d", p.EndTime()),
			"-c:a", "copy",
			audioPath,
		}
		cmd := exec.Command("ffmpeg", argSlice...)
		_, err = cmd.CombinedOutput()
		if err != nil {
			message = err.Error()
			status = Error
			return
		}
		f, err := os.Create(cachedPath)
		if err != nil {
			message = err.Error()
			status = Error
			return
		}
		_, _ = f.Write([]byte(p.FilePath()))
		_ = f.Close()
		message = fmt.Sprintf("Successfully encoded.\nTime taken: %s", dur.String())
		p.encodingState.setState(status, message)
	} else {
		message = fmt.Sprintf("File exist, using cached path %s", cachedPath)
	}
	status = Success
	return
}
