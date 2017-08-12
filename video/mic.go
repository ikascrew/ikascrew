package video

import (
	"math"

	"github.com/gordonklaus/portaudio"
	"github.com/ikascrew/go-opencv/opencv"

	"github.com/ikascrew/ikascrew"
)

func init() {
}

type Microphone struct {
	current int
	frames  []byte
	stream  *portaudio.Stream
	dst     *opencv.IplImage
	bg      *opencv.IplImage
	w       int
	h       int
}

func NewMicrophone() (*Microphone, error) {

	m := Microphone{
		current: 0,
		frames:  make([]byte, 512),
	}

	m.w = ikascrew.Config.Width
	m.h = ikascrew.Config.Height

	m.dst = opencv.CreateImage(m.w, m.h, opencv.IPL_DEPTH_8U, 3)
	m.bg = opencv.CreateImage(m.w, m.h, opencv.IPL_DEPTH_8U, 3)

	//p1 := opencv.Point{0, 288}
	//p2 := opencv.Point{1024, 288}
	//color := opencv.NewScalar(255, 255, 204, 0)
	//opencv.Line(m.bg, p1, p2, color, 5, 8, 0)

	return &m, nil
}

func (m *Microphone) initialize() error {

	inputChannels := 1
	outputChannels := 0
	sampleRate := 8000
	err := portaudio.Initialize()
	if err != nil {
		return err
	}

	stream, err := portaudio.OpenDefaultStream(inputChannels, outputChannels, float64(sampleRate), len(m.frames), m.frames)
	if err != nil {
		return err
	}

	m.stream = stream

	err = m.stream.Start()
	if err != nil {
		return err
	}

	return nil
}

func (m *Microphone) Next() (*opencv.IplImage, error) {

	if m.stream == nil {
		err := m.initialize()
		if err != nil {
			return nil, err
		}

	}

	opencv.Copy(m.bg, m.dst, nil)

	err := m.stream.Read()
	if err != nil {
		return nil, err
	}

	wav := m.frames[:]

	for idx, byt := range wav {

		x := m.w / 2
		y := m.h / 2
		p1 := opencv.Point{x, y}

		t := float64(idx) * 0.5
		var dx float64
		var dy float64

		switch int(t) {
		case 0:
			dx = float64(byt)
			dy = 0.0
		case 1:
			dx = 0.0
			dy = float64(byt)
		case 2:
			dx = float64(byt) * -1.0
			dy = 0.0
		case 3:
			dx = 0.0
			dy = float64(byt) * -1.0
		default:
			dx = math.Cos(t) * float64(byt)
			dy = math.Sin(t) * float64(byt)
		}

		p2 := opencv.Point{x + int(dx)*2, y + int(dy)*2}

		color := opencv.NewScalar(30, 10, 10, 0)
		opencv.Line(m.dst, p1, p2, color, 10, 8, 0)
	}

	m.current++
	if m.current == m.Size() {
		m.current = 0
	}
	return m.dst, nil
}

func (v *Microphone) Wait() int {
	return 33
}

func (v *Microphone) Set(f int) {
	v.current = f
}

func (v *Microphone) Current() int {
	return v.current
}

func (v *Microphone) Size() int {
	return 100
}

func (v *Microphone) Source() string {
	return "ikascrew_Microphone"
}

func (m *Microphone) Release() error {
	m.stream.Stop()
	portaudio.Terminate()

	m.bg.Release()
	m.dst.Release()
	return nil
}
