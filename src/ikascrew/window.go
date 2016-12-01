package ikascrew

import (
	"fmt"
	"github.com/secondarykey/go-opencv/opencv"
	//"time"
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

type MainWindow struct {
	ParentWindow
}

type SubWindow struct {
	stop bool
	ParentWindow
}

func PlaySubWindow(w, v string) error {

	_, ok := windows[w]
	if ok {
		return fmt.Errorf("Exist window name[%s]", w)
	}

	rtn := &SubWindow{}
	win := opencv.NewWindow(w)

	rtn.ParentWindow.window = win
	rtn.stop = false

	win.SetMouseCallback(func(event, x, y, flags int) {
		if flags&opencv.CV_EVENT_LBUTTONDOWN != 0 {
			rtn.stop = !rtn.stop
		}
	})

	q, err := GetVideo(v)
	if err != nil {
		return err
	}

	win.CreateTrackbar("Seek", 1, q.Size(), func(pos int) {
		if pos != q.Current() {
			video, ok := q.(*Video)
			if ok {
				video.cap.SetProperty(opencv.CV_CAP_PROP_POS_FRAMES, float64(pos))
			}
		}
	})

	windows[w] = rtn
	go func() {
		rtn.Play(q)
	}()

	return nil
}

func NewMainWindow(name string) Window {
	rtn := &MainWindow{}
	win := opencv.NewWindow(name)

	rtn.ParentWindow.window = win
	return rtn
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

func (w *MainWindow) Play(q Queue) {

	w.q = q
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

	//t := time.NewTicker(time.Duration(q.Wait()) * time.Millisecond)
	w.q = q
	for {
		//select {
		//case <-t.C:
		if !w.stop {
			img := q.Next()
			if img == nil {
				return
			}
			w.View(img)
			w.ParentWindow.window.SetTrackbarPos("Seek", q.Current())
			opencv.WaitKey(q.Wait())
		}
		//}
	}

}
