package server

import (
	"fmt"
	"time"

	"github.com/golang/glog"
	//"github.com/ikascrew/go-opencv/opencv"

	"github.com/ikascrew/ikascrew"
	pm "github.com/ikascrew/powermate"

	"gocv.io/x/gocv"
)

type Stream struct {
	now_video ikascrew.Video
	now_value float64
	now_image gocv.Mat

	old_video ikascrew.Video
	old_value float64
	old_image gocv.Mat

	release_video ikascrew.Video

	used map[string]bool

	nextFlag bool
	prevFlag bool

	light       float64
	empty_image gocv.Mat
	real_image  gocv.Mat

	wait float64

	mode int
}

const SWITCH_VALUE = 200

const (
	SWITCH = 1
	LIGHT  = 2
	WAIT   = 3
)

func NewStream() (*Stream, error) {
	rtn := Stream{}

	rtn.now_value = 0
	rtn.old_value = 0

	rtn.now_video = nil
	rtn.old_video = nil
	rtn.release_video = nil

	rtn.used = make(map[string]bool)

	rtn.now_image = gocv.NewMatWithSize(ikascrew.Config.Height, ikascrew.Config.Width, gocv.MatTypeCV8UC3)
	rtn.old_image = gocv.NewMatWithSize(ikascrew.Config.Height, ikascrew.Config.Width, gocv.MatTypeCV8UC3)

	rtn.nextFlag = false
	rtn.prevFlag = false

	rtn.empty_image = gocv.NewMatWithSize(ikascrew.Config.Height, ikascrew.Config.Width, gocv.MatTypeCV8UC3)
	rtn.real_image = gocv.NewMatWithSize(ikascrew.Config.Height, ikascrew.Config.Width, gocv.MatTypeCV8UC3)

	rtn.light = 0

	rtn.wait = 0
	rtn.mode = SWITCH
	return &rtn, nil
}

func (s *Stream) Switch(v ikascrew.Video) error {
	if s.used[v.Source()] {
		return fmt.Errorf("until used video")
	}
	s.used[v.Source()] = true

	s.old_value = s.now_value
	s.now_value = 0

	wk := s.release_video
	if wk != nil {
		delete(s.used, wk.Source())
		defer wk.Release()
	}

	s.release_video = s.old_video
	s.old_video = s.now_video
	s.now_video = v

	return nil
}

func (s *Stream) Add(org gocv.Mat) *gocv.Mat {

	alpha := s.light / 200
	gocv.AddWeighted(s.empty_image, float64(alpha), org, float64(1.0-alpha), 0.0, &s.real_image)

	return &s.real_image
}

func (s *Stream) Get(pm bool) (*gocv.Mat, error) {

	old, err := s.getOldImage()
	if err != nil {
		return nil, err
	}

	if old == nil {
		glog.Info("old == nil")
		return s.now_video.Next()
	}

	if !pm {
		if s.now_value != SWITCH_VALUE {
			s.now_value++
		}
	}

	if s.nextFlag {
		if s.now_value == SWITCH_VALUE {
			s.nextFlag = false
		} else if s.now_value < SWITCH_VALUE {
			s.now_value++
		} else {
			s.now_value--
		}
	}

	if s.prevFlag {
		if s.now_value == 0 {
			s.prevFlag = false
		} else if s.now_value > 0 {
			s.now_value--
		} else {
			s.now_value++
		}
	}

	alpha := s.now_value / SWITCH_VALUE

	next, err := s.now_video.Next()
	if err != nil {
		glog.Error("Next video error", err)
		return nil, err
	}

	gocv.AddWeighted(*next, float64(alpha), *old, float64(1.0-alpha), 0.0, &s.now_image)

	return &s.now_image, nil
}

func (s *Stream) getOldImage() (*gocv.Mat, error) {

	if s.release_video == nil {
		if s.old_video != nil {
			return s.old_video.Next()
		}
		return nil, nil
	}

	alpha := s.old_value / SWITCH_VALUE

	next, _ := s.old_video.Next()
	now, _ := s.release_video.Next()

	gocv.AddWeighted(*next, float64(alpha), *now, float64(1.0-alpha), 0.0, &s.old_image)

	return &s.old_image, nil
}

func (s *Stream) Release() {

	s.now_image.Close()
	s.old_image.Close()

	if s.now_video != nil {
		s.now_video.Release()
	}
	if s.old_video != nil {
		s.old_video.Release()
	}
	if s.release_video != nil {
		s.release_video.Release()
	}
}

func (s *Stream) Wait() time.Duration {
	return time.Duration(s.wait + 33.0)
}

func (s *Stream) Effect(e pm.Event) error {
	switch e.Type {
	case pm.Press:
		switch e.Value {
		case pm.Up:
			fmt.Println("Up")
		case pm.Down:
			if s.mode == SWITCH {
				s.mode = LIGHT
				fmt.Println("Light Mode")
			} else if s.mode == LIGHT {
				s.mode = WAIT
				fmt.Println("Wait Mode")
			} else {
				s.mode = SWITCH
				fmt.Println("Switch Mode")
			}
		}
	default:
	}

	switch s.mode {
	case LIGHT:
		switch e.Type {
		case pm.Rotation:
			switch e.Value {
			case pm.Left:
				s.light = s.light + 0.5
			case pm.Right:
				s.light = s.light - 0.5
			}
		}
		fmt.Printf("Light[%f]\n", s.light)
	case SWITCH:
		switch e.Type {
		case pm.Rotation:
			switch e.Value {
			case pm.Left:
				s.now_value--
			case pm.Right:
				s.now_value++
			}
		}
		fmt.Printf("Switch[%f/%d]\n", s.now_value, SWITCH_VALUE)
	case WAIT:
		switch e.Type {
		case pm.Rotation:
			switch e.Value {
			case pm.Left:
				s.wait = s.wait + 0.1
			case pm.Right:
				s.wait = s.wait - 0.1
			}
		}
		fmt.Printf("Wait[%f]\n", s.wait)
	}
	return nil
}

func (s *Stream) PrintVideos(line string) {
	glog.Info(line + "-------------------------------------------------")
	if s.now_video != nil {
		glog.Info("[1]" + s.now_video.Source())
	}

	if s.old_video != nil {
		glog.Info("[2]" + s.old_video.Source())
	}

	if s.release_video != nil {
		glog.Info("[3]" + s.release_video.Source())
	}
}

func (s *Stream) SetSwitch(t string) error {
	if t == "next" {
		s.nextFlag = true
	} else if t == "prev" {
		s.prevFlag = true
	} else {
		return fmt.Errorf("Unknown type[%s]", t)
	}
	return nil
}
