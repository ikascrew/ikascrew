package client

import (
	"C"
	"unsafe"

	"image"
	"image/color"

	"golang.org/x/exp/shiny/screen"

	"github.com/ikascrew/go-opencv/opencv"
	"github.com/ikascrew/ikascrew"
)

type Player struct {
	idx    int
	target []*opencv.IplImage
	*Part
}

func NewPlayer(w screen.Window, s screen.Screen) (*Player, error) {
	p := &Player{}

	r := image.Rect(512, 144, 1536, 720)
	p.Part = &Part{}
	p.Init(w, s, r)

	p.idx = 0
	p.target = make([]*opencv.IplImage, 0)

	return p, nil
}

func (p *Player) setFile(n string) {

	p.idx = 0
	if len(p.target) != 0 {
		work := p.target
		go func() {
			for _, elm := range work {
				elm.Release()
			}
		}()
	}

	p.target = make([]*opencv.IplImage, 5)

	d := ikascrew.ProjectName()
	cap := opencv.NewFileCapture(d + n)
	if cap == nil {
		return
	}
	frames := int(cap.GetProperty(opencv.CV_CAP_PROP_FRAME_COUNT))

	for idx := 0; idx < 5; idx++ {
		frame := int(idx / 5 * int(frames))
		cap.SetProperty(opencv.CV_CAP_PROP_POS_FRAMES, float64(frame))
		img := cap.QueryFrame()
		p.target[idx] = img.Clone()
	}

	cap.Release()

	p.Draw()
}

func (p *Player) Draw() {

	m := p.Part.buffer.RGBA()
	p.idx++
	if len(p.target) == p.idx {
		p.idx = 0
	}
	ipl := p.target[p.idx]

	var height, channels, step int = ipl.Height(), ipl.Channels(), ipl.WidthStep()
	var limg_ptr unsafe.Pointer = ipl.ImageData()
	var data []C.char = (*[1 << 30]C.char)(limg_ptr)[:height*step : height*step]

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
