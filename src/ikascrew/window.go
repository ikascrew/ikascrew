package ikascrew

import (
	"fmt"
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
	window *opencv.Window
}

type MainWindow struct {
	ParentWindow
}

type SubWindow struct {
	ParentWindow
}

func NewSubWindow(name string) Window {

	rtn := &SubWindow{}
	win := opencv.NewWindow("VideoPlayer:" + name)
	stop := false

	win.SetMouseCallback(func(event, x, y, flags int) {
		if flags&opencv.CV_EVENT_LBUTTONDOWN != 0 {
			stop = !stop
			if stop {
				fmt.Printf("status: stop")
			} else {
				fmt.Printf("status: palying")
			}
		}
	})

	rtn.ParentWindow.window = win
	return rtn
}

func NewMainWindow(name string) Window {
	rtn := &MainWindow{}
	win := opencv.NewWindow("VideoPlayer:" + name)
	rtn.ParentWindow.window = win
	return rtn
}

func (w *ParentWindow) View(img *opencv.IplImage) {
	w.window.ShowImage(img)
}

func (w *ParentWindow) Destroy() {
	w.window.Destroy()
}

func (w *MainWindow) Play(q Queue) {
	for {
		img := q.Next()
		if img == nil {
			break
		}
		w.View(img)
		opencv.WaitKey(q.Wait())
	}
}

func (w *SubWindow) Play(q Queue) {

	w.ParentWindow.window.CreateTrackbar("Seek", 1, q.Size(), func(pos int) {
		if pos != q.Current() {
			//cap.SetProperty(opencv.CV_CAP_PROP_POS_FRAMES, float64(pos))
		}
	})

	for {

		img := q.Next()
		if img == nil {
			break
		}

		w.View(img)
		//w.ParentWindow.window.SetTrackberPos("Seek", q.Current())

		key := opencv.WaitKey(q.Wait())
		if key == 27 {
		}

	}
}
