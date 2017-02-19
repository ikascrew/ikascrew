package video

import (
	"github.com/secondarykey/go-opencv/opencv"
)

func init() {
}

type Image struct {
	bg *opencv.IplImage
}

func NewImage() (*Image, error) {
	bg = opencv.LoadImage("projects/20170213/utopia.jpg")
	img := Image{
		bg: bg,
	}
	return &img, nil
}

func (p *Image) Next() *opencv.IplImage {
	return p.bg
}

func (v *Image) Wait() int {
	return 33
}

func (v *Image) Size() int {
	return 100
}

func (v *Image) Current() int {
	return 30
}

func (v *Image) Set(f int) {
}

func (v *Image) Reload() {
}

func (v *Image) Release() {
	v.bg.Release()
}

func (v *Image) Source() string {
	return "_ikascrew_Image.mp4"
}
