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

func NewWindow(name string) *Window {
	rtn := &Window{}
	win := opencv.NewWindow(name)
	rtn.window = win
	return rtn
}

func (w *Window) Play(q *Queue) {

	w.q = q
	for {

		img := q.Next()
		if img != nil {
			w.window.ShowImage(img)
			opencv.WaitKey(q.Wait())
		}
	}
}

func (w *Window) Destroy() {
	if w.q != nil {
		w.q.Release()
	}
	w.window.Destroy()
}
