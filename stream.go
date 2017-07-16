package ikascrew

import (
	"fmt"

	"github.com/secondarykey/go-opencv/opencv"

	pm "github.com/secondarykey/powermate"
)

type Stream struct {
	resource map[string]Video

	now_video Video
	now_value float64
	now_image *opencv.IplImage

	old_video Video
	old_value float64
	old_image *opencv.IplImage

	release_video Video
}

const SWITCH_VALUE = 200

func NewStream() (*Stream, error) {

	s := Stream{}
	s.resource = make(map[string]Video)

	s.now_image = opencv.CreateImage(Config.Width, Config.Height, opencv.IPL_DEPTH_8U, 3)
	s.old_image = opencv.CreateImage(Config.Width, Config.Height, opencv.IPL_DEPTH_8U, 3)

	s.now_value = 0
	s.old_value = 0

	s.now_video = nil
	s.old_video = nil
	s.release_video = nil

	return &s, nil
}

func (s *Stream) Push(v Video, support bool) error {

	if s.release_video != nil {
		return fmt.Errorf("Until Switch")
	}

	_, ok := s.resource[v.Source()]
	if ok {
		return fmt.Errorf("Exist Video[%s]", v.Source())
	}

	s.resource[v.Source()] = v

	//次が来たらMateの対象にする
	s.old_value = s.now_value
	s.now_value = 0

	s.release_video = s.old_video
	s.old_video = s.now_video
	s.now_video = v

	fmt.Println("Push")

	return nil
}

func (s *Stream) Next(sw bool) (*opencv.IplImage, error) {

	old, err := s.getOldImage()
	if err != nil {
		return nil, err
	}

	if old == nil {
		return s.now_video.Next()
	}

	if !sw {
		s.now_value++
		// 切り替え前のビデオを削除
		if s.now_value == SWITCH_VALUE {
			defer func() {
				s.old_video.Release()
				s.old_video = nil
			}()
		}
	}

	alpha := s.now_value / SWITCH_VALUE

	next, _ := s.now_video.Next()
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

	//完全に切り替える方向に持っていく
	if s.old_value > SWITCH_VALUE {
		s.old_value--
	} else if s.old_value < SWITCH_VALUE {
		s.old_value++
	}

	alpha := s.old_value / SWITCH_VALUE

	next, _ := s.old_video.Next()
	now, _ := s.release_video.Next()

	opencv.AddWeighted(next, float64(alpha), now, float64(1.0-alpha), 0.0, s.old_image)

	if s.old_value == SWITCH_VALUE {
		s.release_video.Release()

		s.release_video = nil
	}
	return s.old_image, nil
}

func (s *Stream) Wait() int {
	return s.now_video.Wait()
}

func (s *Stream) Release() error {

	//Stream のリリースは終了時のみ行う
	return nil

}

func (s *Stream) Effect(e pm.Event) error {

	switch e.Type {
	case pm.Rotation:
		switch e.Value {
		case pm.Left:
			s.new_value--
		case pm.Right:
			s.new_value++
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
