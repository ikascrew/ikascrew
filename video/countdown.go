package video

import (
	"fmt"
	"image"
	"image/color"
	"time"

	"github.com/ikascrew/ikascrew"
	"gocv.io/x/gocv"
)

func init() {
}

var loc, _ = time.LoadLocation("Asia/Tokyo")

var Target = time.Date(2020, time.January, 1, 0, 0, 0, 0, loc)

//var Target = time.Date(2019, time.December, 31, 9, 2, 0, 0, loc)
var jst = time.FixedZone("Asia/Tokyo", 9*60*60)

const Line1 = "Happy"
const Line2 = "New Rising"

type Countdown struct {
	frames int
	name   string
	source *gocv.Mat

	counter int
	cap     *gocv.VideoCapture

	target int64
}

func NewCountdown(file string) (*Countdown, error) {

	f := Countdown{
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
	f.target = Target.In(jst).Unix()
	return &f, nil
}

func (v *Countdown) Next() (*gocv.Mat, error) {

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

	now := time.Now().In(jst)
	d := v.target - now.Unix()

	if d >= 0 {

		buf := fmt.Sprintf("%d", d)

		// 3 200
		// 2 300
		// 1 400

		left := 500 - (len(buf) * 100)
		gocv.PutText(v.source, buf, image.Pt(left, 400),
			gocv.FontHersheyComplexSmall, 16.0, color.RGBA{255, 255, 255, 0}, 4)

		if len(buf) <= 1 {
			//gocv.Ellipse(v.source, image.Pt(512, 288), image.Pt(712, 288), 0, 0, 3.14, color.RGBA{255, 255, 255, 0}, 2)
			gocv.Circle(v.source, image.Pt(502, 295), 200, color.RGBA{255, 255, 255, 0}, 8)
		}

	} else {
		gocv.PutText(v.source, Line1, image.Pt(180, 200),
			gocv.FontHersheyComplexSmall, 9.0, color.RGBA{255, 255, 255, 0}, 4)
		gocv.PutText(v.source, Line2, image.Pt(10, 450),
			gocv.FontHersheyComplexSmall, 7.4, color.RGBA{255, 255, 255, 0}, 4)

	}

	return v.source, nil
}

func (v *Countdown) Set(f int) {
	v.cap.Set(gocv.VideoCapturePosFrames, float64(f))
}

func (v *Countdown) Current() int {
	return int(v.cap.Get(gocv.VideoCapturePosFrames))
}

func (v *Countdown) Size() int {
	return v.frames
}

func (v *Countdown) Source() string {
	return v.name
}

func (v *Countdown) Release() error {
	if v.cap != nil {
		v.cap.Close()
	}
	v.cap = nil
	return nil
}
