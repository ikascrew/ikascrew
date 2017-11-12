package server

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/ikascrew/go-opencv/opencv"

	"github.com/ikascrew/ikascrew"
	pm "github.com/ikascrew/powermate"
)

func init() {
}

type Window struct {
	name string
	wait chan ikascrew.Video

	stream *Stream

	PowerMate bool
}

func NewWindow(name string) (*Window, error) {

	rtn := &Window{}

	rtn.name = name
	rtn.wait = make(chan ikascrew.Video)

	var err error
	rtn.stream, err = NewStream()
	return rtn, err
}

func (w *Window) Push(v ikascrew.Video) error {
	w.stream.PrintVideos("Push Start")
	w.wait <- v
	w.stream.PrintVideos("Push End")
	return nil
}

func (w *Window) Play(v ikascrew.Video) error {

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	opencv.StartWindowThread()
	win := opencv.NewWindow(w.name)
	defer win.Destroy()

	err := w.stream.Switch(v)
	if err != nil {
		return err
	}

	for {
		glog.Info("Main Loop")
		select {
		case v := <-w.wait:
			err := w.stream.Switch(v)
			if err != nil {
				glog.Error("Stream Push Error:", err)
			}
		default:
			err := w.Display(win)
			if err != nil {
				glog.Error("Window Display Error:", err)
			}
		}
		glog.Info("Main Loop End")
	}

	return fmt.Errorf("Error : Stream is nil")
}

func (w *Window) Display(win *opencv.Window) error {

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		//fps 30
		time.Sleep(w.stream.Wait() * time.Millisecond)
		wg.Done()
	}()

	img, err := w.stream.Get(w.PowerMate)
	if err != nil {
		return err
	}
	win.ShowImage(w.stream.Add(img))

	wg.Wait()

	return nil
}

func (w *Window) SetSwitch(t string) error {
	return w.stream.SetSwitch(t)
}

func (w *Window) Effect(e pm.Event) error {
	return w.stream.Effect(e)
}

func (w *Window) Destroy() {
	w.stream.Release()
}

func (w *Window) FullScreen() {
	//w.window.SetProperty(opencv.CV_WND_PROP_FULLSCREEN, float64(opencv.CV_WINDOW_FULLSCREEN))
}
