package video

import (
	"fmt"

	"github.com/secondarykey/go-opencv/opencv"
)

func init() {
}

type Image struct {
	current int
	name    string
	bg      *opencv.IplImage
}

func NewImage(f string) (*Image, error) {
	img := Image{
		name:    f,
		current: 0,
	}
	img.bg = opencv.LoadImage(f)
	if img.bg == nil {
		return nil, fmt.Errorf("Error:LoadImage[%s]", f)
	}

	return &img, nil
}

func (v *Image) Next() (*opencv.IplImage, error) {
	v.current++
	if v.current == v.Size() {
		v.current = 0
	}
	return v.bg, nil
}

func (v *Image) Wait() int {
	return 33
}

func (v *Image) Set(f int) {
	v.current = f
}

func (v *Image) Current() int {
	return v.current
}

func (v *Image) Size() int {
	return 100
}

func (v *Image) Source() string {
	return v.name
}

func (v *Image) Release() error {
	v.bg.Release()
	return nil
}
