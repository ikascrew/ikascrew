package ikascrew

import (
	"fmt"
	"github.com/secondarykey/go-opencv/opencv"
	pm "github.com/secondarykey/powermate"
)

func init() {
}

type Window struct {
	stream    *Stream
	window    *opencv.Window
	PowerMate bool
}

func NewWindow(name string) *Window {
	rtn := &Window{}
	win := opencv.NewWindow(name)
	rtn.window = win
	rtn.stream = &Stream{}
	rtn.PowerMate = false
	return rtn
}

func (w *Window) Play(v Video) error {
	w.stream.Push(v)
	for {
		img, err := w.stream.Next()
		if err != nil {
			return err
		}

		if img != nil {
			w.window.ShowImage(img)
			opencv.WaitKey(w.stream.Wait())
		} else {
			fmt.Println("Next() Image Nil!!!")
		}
	}
	return fmt.Errorf("Error : Stream is nil")
}

func (w *Window) Event(e pm.Event) error {
	switch e.Type {
	case pm.Rotation:
		switch e.Value {
		case om.Left:
		case om.Right:
		}
	case pm.Press:
		switch e.Value {
		case om.Up:
		case om.Down:
		}
	default:
	}
}

func (w *Window) Destroy() {
	w.stream.Release()
	w.window.Destroy()
}

func (w *Window) FullScreen() {
	w.window.SetProperty(opencv.CV_WND_PROP_FULLSCREEN, float64(opencv.CV_WINDOW_FULLSCREEN))
}
