package effect

import (
	"fmt"

	"github.com/secondarykey/go-opencv/opencv"
	"github.com/secondarykey/ikascrew"
)

type Normal struct {
	video ikascrew.Video
}

func NewNormal(v ikascrew.Video) (*Normal, error) {
	l := Normal{
		video: v,
	}
	return &l, nil
}

func (e *Normal) Next() (*opencv.IplImage, error) {
	return e.video.Next()
}

func (e *Normal) Wait() int {
	return e.video.Wait()
}

func (e *Normal) Release() error {
	return e.video.Release()
}

func (e *Normal) String() string {
	return fmt.Sprintf(`{
		"source" : "%s",
		"current" : %d
	}`, e.video.Source(), e.video.Current())
}
