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

func NewWindow(name string) (*Window, error) {
	var err error
	rtn := &Window{}
	win := opencv.NewWindow(name)
	rtn.window = win
	rtn.stream, err = NewStream()
	if err != nil {
		return nil, err
	}
	rtn.PowerMate = false
	return rtn, nil
}

func (w *Window) Push(v Video) error {
	return w.stream.Push(v)
}

func (w *Window) Play(v Video) error {

	err := w.Push(v)
	if err != nil {
		return err
	}

	for {
		img, err := w.stream.Next(w.PowerMate)
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

func (w *Window) Effect(e pm.Event) error {
	return w.stream.Effect(e)
}

func (w *Window) Destroy() {
	w.stream.Release()
	w.window.Destroy()
}

func (w *Window) FullScreen() {
	w.window.SetProperty(opencv.CV_WND_PROP_FULLSCREEN, float64(opencv.CV_WINDOW_FULLSCREEN))
}

func (w *Window) Now() {
	//TODO 現状の表示を取得
}
