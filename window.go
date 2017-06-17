package ikascrew

import (
	"fmt"
	"github.com/secondarykey/go-opencv/opencv"
)

func init() {
}

type Window struct {
	effect  Effect
	old     Effect
	release Effect
	window  *opencv.Window
}

func NewWindow(name string, e Effect) *Window {

	rtn := &Window{}
	win := opencv.NewWindow(name)
	rtn.window = win
	rtn.effect = e
	rtn.old = nil
	rtn.release = nil

	return rtn
}

func (w *Window) Play() error {
	for {

		img, err := w.effect.Next()
		if err != nil {
			return err
		}

		if img != nil {
			w.window.ShowImage(img)
			opencv.WaitKey(w.effect.Wait())
		}
	}
	return fmt.Errorf("Error : Stream is nil")
}

func (w *Window) Destroy() {
	if w.effect != nil {
		w.effect.Release()
	}
	w.window.Destroy()
}

func (w *Window) Current() string {
	return w.effect.String()
}

func (w *Window) FullScreen() {
	w.window.SetProperty(opencv.CV_WND_PROP_FULLSCREEN, float64(opencv.CV_WINDOW_FULLSCREEN))
}

func (w *Window) GetEffect() Effect {
	return w.effect
}

func (w *Window) SetEffect(e Effect) {
	if w.release != nil {
		w.release.Release()
	}
	if w.old != nil {
		w.release = w.old
	}
	w.old = w.effect

	w.effect = e
}
