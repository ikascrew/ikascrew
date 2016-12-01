package ikascrew

import (
	"fmt"
	"github.com/secondarykey/go-opencv/opencv"
)

func init() {
}

type Video struct {
	FPS      int
	Frames   int
	Position int
	cap      *opencv.Capture
	file     string
}

func NewVideo(f string) (*Video, error) {

	cap := opencv.NewFileCapture(f)
	if cap == nil {
		return nil, fmt.Errorf("New Capture Error:[%s]", f)
	}

	fps := int(cap.GetProperty(opencv.CV_CAP_PROP_FPS))
	frames := int(cap.GetProperty(opencv.CV_CAP_PROP_FRAME_COUNT))
	v := &Video{
		FPS:    fps,
		Frames: frames,
		cap:    cap,
		file:   f,
	}

	return v, nil
}

func (v *Video) Next() *opencv.IplImage {

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

func (v *Video) Wait() int {
	return 1000 / v.FPS
}

func (v *Video) Size() int {
	return v.Frames
}

func (v *Video) Current() int {
	return v.Position
}

func (v *Video) Reload() {
	v.cap.SetProperty(opencv.CV_CAP_PROP_POS_FRAMES, float64(0))
}

func (v *Video) Release() {
	cp := v.cap
	v.cap = nil

	cp.Release()
}
