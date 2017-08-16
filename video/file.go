package video

import (
	"fmt"

	"github.com/ikascrew/go-opencv/opencv"
)

func init() {
}

type File struct {
	frames int
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

	f.frames = int(f.cap.GetProperty(opencv.CV_CAP_PROP_FRAME_COUNT))

	return &f, nil
}

func (v *File) Next() (*opencv.IplImage, error) {

	if v.cap == nil {
		return nil, fmt.Errorf("Error:Caputure is nil")
	}

	pos := int(v.cap.GetProperty(opencv.CV_CAP_PROP_POS_FRAMES))
	if pos == v.Size() {
		v.Set(1)
	}

	img := v.cap.QueryFrame()
	if img == nil {
		v.Set(1)
		return nil, fmt.Errorf("Error:Image is nil")
	}

	return img, nil
}

func (v *File) Set(f int) {
	v.cap.SetProperty(opencv.CV_CAP_PROP_POS_FRAMES, float64(f))
}

func (v *File) Current() int {
	return int(v.cap.GetProperty(opencv.CV_CAP_PROP_POS_FRAMES))
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
