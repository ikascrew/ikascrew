package effect

import (
	"github.com/secondarykey/go-opencv/opencv"
	"github.com/secondarykey/ikascrew"
)

type Switch struct {
	video ikascrew.Video

	now    ikascrew.Effect
	count  float64
	number float64
	img    *opencv.IplImage
}

func NewSwitch(v ikascrew.Video, e ikascrew.Effect) (*Switch, error) {

	//変更用の画像を確保
	//TODO 設定値による変更
	dst := opencv.CreateImage(1024, 576, opencv.IPL_DEPTH_8U, 3)

	s := Switch{
		video:  v,
		now:    e,
		count:  1,
		number: 200,
		img:    dst,
	}
	return &s, nil
}

func (s *Switch) Next() (*opencv.IplImage, error) {

	if s.count == s.number {

		if s.img != nil {
			defer s.finish()
		}

		return s.video.Next()
	} else {

		now, _ := s.now.Next()
		next, _ := s.video.Next()

		alpha := s.count / s.number
		opencv.AddWeighted(next, float64(alpha), now, float64(1.0-alpha), 0.0, s.img)

		s.count++

		return s.img, nil
	}
}

func (s *Switch) finish() error {
	s.img.Release()
	s.now.Release()

	s.img = nil
	s.now = nil
	return nil
}

func (s *Switch) Wait() int {
	if s.count == s.number {
		return s.video.Wait()
	}
	return s.now.Wait()
}

func (s *Switch) Release() error {
	return s.video.Release()
}

func (s *Switch) String() string {
	return ""
}
