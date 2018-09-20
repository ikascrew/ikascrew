package video

import (
	"fmt"

	"github.com/ikascrew/ikascrew"
	"gocv.io/x/gocv"
)

func init() {
}

type File struct {
	frames int
	name   string
	source *gocv.Mat

	cap *gocv.VideoCapture
}

func NewFile(file string) (*File, error) {

	f := File{
		name: file,
	}
	var err error

	f.cap, err = gocv.VideoCaptureFile(file)
	if err != nil {
		return nil, err
	}

	if f.cap == nil {
		return nil, fmt.Errorf("New Capture Error:[%s]", f)
	}

	//f.frames = int(f.cap.GetProperty(opencv.CV_CAP_PROP_FRAME_COUNT))
	f.frames = int(f.cap.Get(gocv.VideoCaptureFrameCount))
	v := gocv.NewMatWithSize(ikascrew.Config.Height, ikascrew.Config.Width, gocv.MatTypeCV8UC3)

	f.source = &v
	return &f, nil
}

func (v *File) Next() (*gocv.Mat, error) {

	if v.cap == nil {
		return nil, fmt.Errorf("Error:Caputure is nil")
	}

	//pos := int(v.cap.GetProperty(opencv.CV_CAP_PROP_POS_FRAMES))
	pos := int(v.cap.Get(gocv.VideoCapturePosFrames))
	if pos == v.Size() {
		v.Set(1)
	}

	v.cap.Read(v.source)
	if v.source.Empty() {
		v.Set(1)
		return nil, fmt.Errorf("Error:Image is nil")
	}

	return v.source, nil
}

func (v *File) Set(f int) {
	v.cap.Set(gocv.VideoCapturePosFrames, float64(f))
}

func (v *File) Current() int {
	return int(v.cap.Get(gocv.VideoCapturePosFrames))
}

func (v *File) Size() int {
	return v.frames
}

func (v *File) Source() string {
	return v.name
}

func (v *File) Release() error {
	if v.cap != nil {
		v.cap.Close()
	}
	v.cap = nil
	return nil
}
