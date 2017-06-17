package effect

import (
	"fmt"

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

	fmt.Printf("Switch[%p] New\n", &s)

	return &s, nil
}

func (s *Switch) Next() (*opencv.IplImage, error) {

	if s.video == nil || s.now == nil {
		fmt.Println("Caution [Video == nil]")
		return nil, nil
	}

	if s.count == s.number {
		return s.video.Next()
	} else {

		now, _ := s.now.Next()
		next, _ := s.video.Next()

		alpha := s.count / s.number
		opencv.AddWeighted(next, float64(alpha), now, float64(1.0-alpha), 0.0, s.img)

		s.count++

		if s.count == s.number {
			fmt.Println("Switch Done!")
		}
		return s.img, nil
	}
}

func (s *Switch) Wait() int {
	if s.count == s.number {
		return s.video.Wait()
	}
	return s.now.Wait()
}

func (s *Switch) Release() error {

	fmt.Printf("Switch[%p] Release\n", s)

	if s.img != nil {
		s.img.Release()
	}
	s.img = nil

	if s.now != nil {
		s.now.Release()
	}
	s.now = nil

	if s.video != nil {
		err := s.video.Release()
		if err != nil {
			return err
		}
	}

	s.video = nil
	return nil
}

func (s *Switch) String() string {
	return ""
}
