package ikascrew

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/ikascrew/go-opencv/opencv"
	pm "github.com/ikascrew/powermate"
)

type Stream struct {
	now_video Video
	now_value float64
	now_image *opencv.IplImage

	old_video Video
	old_value float64
	old_image *opencv.IplImage

	release_video Video

	used map[string]bool

	nextFlag bool
	prevFlag bool
}

const SWITCH_VALUE = 200

func NewStream() (*Stream, error) {
	rtn := Stream{}

	rtn.now_value = 0
	rtn.old_value = 0

	rtn.now_video = nil
	rtn.old_video = nil
	rtn.release_video = nil

	rtn.used = make(map[string]bool)

	rtn.now_image = opencv.CreateImage(Config.Width, Config.Height, opencv.IPL_DEPTH_8U, 3)
	rtn.old_image = opencv.CreateImage(Config.Width, Config.Height, opencv.IPL_DEPTH_8U, 3)

	rtn.nextFlag = false
	rtn.prevFlag = false

	return &rtn, nil
}

func (s *Stream) Switch(v Video) error {
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

func (s *Stream) Get(pm bool) (*opencv.IplImage, error) {

	old, err := s.getOldImage()
	if err != nil {
		fmt.Println(err)
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

	opencv.AddWeighted(next, float64(alpha), old, float64(1.0-alpha), 0.0, s.now_image)

	return s.now_image, nil
}

func (s *Stream) getOldImage() (*opencv.IplImage, error) {

	if s.release_video == nil {
		if s.old_video != nil {
			return s.old_video.Next()
		}
		return nil, nil
	}

	alpha := s.old_value / SWITCH_VALUE

	next, _ := s.old_video.Next()
	now, _ := s.release_video.Next()

	opencv.AddWeighted(next, float64(alpha), now, float64(1.0-alpha), 0.0, s.old_image)

	return s.old_image, nil
}

func (s *Stream) Release() {

	s.now_image.Release()
	s.old_image.Release()

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

func (s *Stream) Effect(e pm.Event) error {
	switch e.Type {
	case pm.Rotation:
		switch e.Value {
		case pm.Left:
			s.now_value--
		case pm.Right:
			s.now_value++
		}
	case pm.Press:
		switch e.Value {
		case pm.Up:
		case pm.Down:
		}
	default:
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
