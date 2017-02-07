package ikascrew

import (
	"fmt"
	"github.com/secondarykey/go-opencv/opencv"
)

func init() {
}

type File struct {
	FPS      int
	Frames   int
	Position int
	cap      *opencv.Capture
	Name     string
}

func NewFile(f string) (*File, error) {

	cap := opencv.NewFileCapture(f)
	if cap == nil {
		return nil, fmt.Errorf("New Capture Error:[%s]", f)
	}

	fps := int(cap.GetProperty(opencv.CV_CAP_PROP_FPS))
	frames := int(cap.GetProperty(opencv.CV_CAP_PROP_FRAME_COUNT))
	v := &File{
		FPS:    fps,
		Frames: frames,
		cap:    cap,
		Name:   f,
	}

	return v, nil
}

func (v *File) Next() *opencv.IplImage {

	if v.cap == nil {
		return nil
	}

	img := v.cap.QueryFrame()
	v.Position = int(v.cap.GetProperty(opencv.CV_CAP_PROP_POS_FRAMES))

	if v.Position >= v.Frames {
		v.Reload()
	}
	return img
}

func (v *File) Wait() int {
	return 1000 / v.FPS
}

func (v *File) Size() int {
	return v.Frames
}

func (v *File) Current() int {
	return v.Position
}

func (v *File) Set(f int) {

	if f > v.Frames {
		f = f % v.Frames
	}

	v.cap.SetProperty(opencv.CV_CAP_PROP_POS_FRAMES, float64(f))
}

func (v *File) Reload() {
	v.Set(0)
}

func (v *File) Release() {
	cp := v.cap
	cp.Release()
}

func (v *File) Source() string {
	return v.Name
}
