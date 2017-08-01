package video

import (
	"fmt"

	"github.com/ikascrew/go-opencv/opencv"
)

func init() {
}

type File struct {
	fps    int
	frames int
	pos    int
	name   string

	cap *opencv.Capture
}

func NewFile(file string) (*File, error) {

	f := File{
		name: file,
	}

	f.cap = opencv.NewFileCapture(file)
	if f.cap == nil {
		return nil, fmt.Errorf("New Capture Error:[%s]", f)
	}

	f.fps = int(f.cap.GetProperty(opencv.CV_CAP_PROP_FPS))
	f.frames = int(f.cap.GetProperty(opencv.CV_CAP_PROP_FRAME_COUNT))

	return &f, nil
}

func (v *File) Next() (*opencv.IplImage, error) {

	if v.cap == nil {
		return nil, fmt.Errorf("Error:Caputure is nil")
	}

	img := v.cap.QueryFrame()
	if img == nil {
		return nil, fmt.Errorf("Error:Image is nil")
	}

	v.pos = int(v.cap.GetProperty(opencv.CV_CAP_PROP_POS_FRAMES))
	if v.pos == v.Size() {
		v.Set(0)
	}

	return img, nil
}

func (v *File) Wait() int {
	return 1000 / v.fps
}

func (v *File) Set(f int) {
	if f > v.frames {
		f = f % v.frames
	}
	v.cap.SetProperty(opencv.CV_CAP_PROP_POS_FRAMES, float64(f))
	v.pos = f
}

func (v *File) Current() int {
	return v.pos
}

func (v *File) Size() int {
	return v.frames
}

func (v *File) Source() string {
	return v.name
}

func (v *File) Release() error {
	if v.cap != nil {
		v.cap.Release()
	}
	v.cap = nil
	return nil
}
