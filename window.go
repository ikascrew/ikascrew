package ikascrew

import (
	"fmt"
	"sync"

	"github.com/ikascrew/go-opencv/opencv"
	pm "github.com/ikascrew/powermate"

	"github.com/golang/glog"
)

func init() {
}

type Window struct {
	stream *Stream
	window *opencv.Window

	m         *sync.Mutex
	wait      chan Video
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

	rtn.m = new(sync.Mutex)
	rtn.wait = make(chan Video)

	return rtn, nil
}

func (w *Window) Push(v Video) error {

	w.stream.PrintVideos("Push Start")
	w.wait <- v
	w.stream.PrintVideos("Push End")
	return nil
}

func (w *Window) Play(v Video) error {

	err := w.stream.Push(v)
	if err != nil {
		return err
	}

	for {
		//w.m.Lock()
		glog.Info("Main Loop")
		select {
		case v := <-w.wait:
			err := w.stream.Push(v)
			if err != nil {
				glog.Error("Stream Push Error:", err)
				//TODO もう一回回るようにする
			}
		default:
			img, err := w.stream.Next(w.PowerMate)
			if err != nil {
				glog.Error("Stream Next Error:", err)
			}

			if img != nil {
				fmt.Println("ShowImage")
				w.window.ShowImage(img)
				if err != nil {
					glog.Error("Window ShowImage Error:", err)
				}
				fmt.Printf("Wait(%d)\n", w.stream.Wait())
				opencv.WaitKey(w.stream.Wait())
			} else {
				glog.Error("Next Image nil")
			}
		}
		glog.Info("Main Loop End")
		//w.m.Unlock()
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
	//w.window.SetProperty(opencv.CV_WND_PROP_FULLSCREEN, float64(opencv.CV_WINDOW_FULLSCREEN))
}
