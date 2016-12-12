package ikascrew

import (
	"github.com/secondarykey/go-opencv/opencv"
)

func init() {
}

type Window interface {
	View(*opencv.IplImage)
	Play(Queue)
	Destroy()
}

type ParentWindow struct {
	q      Queue
	window *opencv.Window
}

func (w *ParentWindow) View(img *opencv.IplImage) {
	w.window.ShowImage(img)
}

func (w *ParentWindow) Destroy() {
	if w.q != nil {
		w.q.Release()
	}
	w.window.Destroy()
}
