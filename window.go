package ikascrew

import (
	"github.com/secondarykey/go-opencv/opencv"
)

func init() {
}

type Window struct {
	q      *Queue
	window *opencv.Window
}

func NewWindow(name string, q *Queue) *Window {
	rtn := &Window{}
	win := opencv.NewWindow(name)
	rtn.window = win

	rtn.q = q

	return rtn
}

func (w *Window) Play() {

	for {
		img := w.q.Next()
		if img != nil {
			w.window.ShowImage(img)
			opencv.WaitKey(w.q.Wait())
		}
	}
}

func (w *Window) Destroy() {
	if w.q != nil {
		w.q.Release()
	}
	w.window.Destroy()
}

func (w *Window) FullScreen() {
	w.window.SetProperty(opencv.CV_WND_PROP_FULLSCREEN, opencv.CV_WINDOW_FULLSCREEN)
}
