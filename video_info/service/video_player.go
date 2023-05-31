package service

import (
	"fmt"
	"time"
)

type Video struct {
	duration  int64
	timeStart time.Time
}

func NewVideo(duration int64) *Video {
	return &Video{duration: duration}
}

func (video *Video) Start() {
	go video.timer()
}

func (video *Video) timer() {
	nanoseconds := video.duration * 1e6
	i := 0

	for {
		i++
		video.timeStart = time.Now()
		timer1 := time.NewTimer(time.Duration(nanoseconds))
		<-timer1.C
		fmt.Printf("Round %d finished", i)
	}
}

func (video *Video) GetCurrentDuration() int64 {
	return int64(time.Since(video.timeStart)) / 1e6
}
