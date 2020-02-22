package video

import (
	"fmt"
	"log"

	"github.com/ikascrew/ikascrew"
	"gocv.io/x/gocv"
)

func init() {
}

type LoopFile struct {
	frames int
	name   string

	rtn    *gocv.Mat
	origin *gocv.Mat
	remix  *gocv.Mat

	originCap *gocv.VideoCapture
	remixCap  *gocv.VideoCapture
}

func NewLoopFile(file string) (*LoopFile, error) {

	f := LoopFile{
		name: file,
	}
	var err error

	f.originCap, err = gocv.VideoCaptureFile(file)
	if err != nil {
		return nil, err
	}

	if f.originCap == nil {
		return nil, fmt.Errorf("New Capture Error:[%s]", f)
	}

	f.remixCap, err = gocv.VideoCaptureFile(file)
	if err != nil {
		return nil, err
	}

	if f.remixCap == nil {
		return nil, fmt.Errorf("New Capture Error:[%s]", f)
	}

	//f.frames = int(f.cap.GetProperty(opencv.CV_CAP_PROP_FRAME_COUNT))
	f.frames = int(f.originCap.Get(gocv.VideoCaptureFrameCount))
	v := gocv.NewMatWithSize(ikascrew.Config.Height, ikascrew.Config.Width, gocv.MatTypeCV8UC3)

	o := gocv.NewMatWithSize(ikascrew.Config.Height, ikascrew.Config.Width, gocv.MatTypeCV8UC3)
	r := gocv.NewMatWithSize(ikascrew.Config.Height, ikascrew.Config.Width, gocv.MatTypeCV8UC3)

	f.rtn = &v

	f.origin = &o
	f.remix = &r
	return &f, nil
}

func (v *LoopFile) Next() (*gocv.Mat, error) {

	if v.originCap == nil {
		return nil, fmt.Errorf("Error:Caputure is nil")
	}

	size := int(float64(v.frames) * 0.2)
	remixPos := v.frames - size

	v.originCap.Read(v.origin)

	pos := int(v.originCap.Get(gocv.VideoCapturePosFrames))
	//位置がまだの場合
	if pos < remixPos {
		log.Printf("Origin=%d \n", pos)
		return v.origin, nil
	} else if pos == remixPos {
		v.remixCap.Set(gocv.VideoCapturePosFrames, float64(1))
	}

	v.remixCap.Read(v.remix)

	alpha := float64(pos-remixPos) / float64(size)
	gocv.AddWeighted(*v.remix, float64(alpha), *v.origin, float64(1.0-alpha), 0.0, v.rtn)

	pos = int(v.originCap.Get(gocv.VideoCapturePosFrames))
	if pos == v.Size() {
		v.Set(1)
	}

	return v.rtn, nil
}

func (v *LoopFile) Set(f int) {
	size := int(float64(v.frames) * 0.2)
	v.originCap.Set(gocv.VideoCapturePosFrames, float64(size+f))
}

func (v *LoopFile) Current() int {
	return int(v.originCap.Get(gocv.VideoCapturePosFrames))
}

func (v *LoopFile) Size() int {
	return v.frames
}

func (v *LoopFile) Source() string {
	return v.name
}

func (v *LoopFile) Release() error {
	if v.originCap != nil {
		v.originCap.Close()
	}

	if v.remixCap != nil {
		v.remixCap.Close()
	}

	v.originCap = nil
	v.remixCap = nil
	return nil
}
