package client

import (
	"image"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
)

type Window struct {
	window screen.Window
	buffer screen.Buffer
}

func NewWindow(t string, w, h int) (window *Window, err error) {

	window = &Window{
		window: nil,
		buffer: nil,
	}

	driver.Main(func(s screen.Screen) {
		opt := &screen.NewWindowOptions{
			Title:  t,
			Width:  w,
			Height: h,
		}
		window.window, err = s.NewWindow(opt)
		if err != nil {
			return
		}
		winSize := image.Point{w, h}
		window.buffer, err = s.NewBuffer(winSize)
		if err != nil {
			return
		}
	})
	return
}

func (w *Window) Release() {
	w.window.Release()
	w.buffer.Release()
}
