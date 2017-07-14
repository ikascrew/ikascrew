package ikascrew

import (
	"github.com/secondarykey/go-opencv/opencv"
)

type Stream struct {
	v Video
}

func (s *Stream) Push(v Video) error {

	s.v = v

	//次が来たらMateの対象にする

	//現在との境界値はそのまま

	//見ているビデオがなくなったら、リリースして、管理から削除

	return nil
}

func (s *Stream) Next() (*opencv.IplImage, error) {
	return s.v.Next()
}

func (s *Stream) Wait() int {
	return s.v.Wait()
}

func (s *Stream) Release() error {

	//Stream のリリースは終了時のみ行う

}
