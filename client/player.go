package client

import (
	"C"

	"image"
	"image/color"

	"golang.org/x/exp/shiny/screen"

	"github.com/ikascrew/ikascrew"

	"gocv.io/x/gocv"
)

type Player struct {
	name   string
	idx    int
	target []gocv.Mat
	*Part
}

func NewPlayer(w screen.Window, s screen.Screen) (*Player, error) {
	p := &Player{}

	r := image.Rect(512, 144, 1536, 720)
	p.Part = &Part{}
	p.Init(w, s, r)

	p.idx = 0
	p.target = make([]gocv.Mat, 0)

	return p, nil
}

func (p *Player) setFile(n string) {

	if p.name == n {
		return
	}

	p.idx = 0
	if len(p.target) != 0 {
		work := p.target
		go func() {
			for _, elm := range work {
				elm.Close()
			}
		}()
	}

	p.name = n
	p.target = make([]gocv.Mat, 5)

	d := ikascrew.ProjectName()
	cap, err := gocv.VideoCaptureFile(d + n)
	if err != nil {
		return
	}

	if cap == nil {
		return
	}
	defer cap.Close()
	frames := cap.Get(gocv.VideoCaptureFrameCount)

	for idx := 0; idx < 5; idx++ {

		frame := int(float64(idx) / 5.0 * float64(frames))
		cap.Set(gocv.VideoCapturePosFrames, float64(frame))

		mat := gocv.NewMat()
		cap.Read(&mat)
		p.target[idx] = mat
	}

}

func (p *Player) Draw() {

	if len(p.target) <= 0 {
		return
	}

	m := p.Part.buffer.RGBA()
	p.idx++
	if len(p.target) == p.idx {
		p.idx = 0
	}
	ipl := p.target[p.idx]

	var height, channels, step int = ipl.Rows(), ipl.Channels(), ipl.Step()

	data := ipl.DataPtrUint8()
	c := color.NRGBA{R: uint8(0), G: uint8(0), B: uint8(0), A: uint8(255)}
	for y := 0; y < height; y++ {
		for x := 0; x < step; x = x + channels {
			c.B = uint8(data[y*step+x])
			c.G = uint8(data[y*step+x+1])
			c.R = uint8(data[y*step+x+2])
			if channels == 4 {
				c.A = uint8(data[y*step+x+3])
			}
			m.Set(int(x/channels), y, c)
		}
	}
	return
}
