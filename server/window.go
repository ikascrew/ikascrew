package server

import (
	"fmt"
	"image"
	"image/color"
	"runtime"
	"sync"
	"time"

	"github.com/golang/glog"

	"github.com/ikascrew/ikascrew"
	pm "github.com/ikascrew/powermate"

	"gocv.io/x/gocv"
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

	win := gocv.NewWindow(w.name)
	defer win.Close()

	win.MoveWindow(0, 0)
	win.ResizeWindow(1024, 576)

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

var counter = 0

func (w *Window) Display(win *gocv.Window) error {

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		time.Sleep(w.stream.Wait() * time.Millisecond)
		wg.Done()
	}()

	img, err := w.stream.Get(w.PowerMate)
	if err != nil {
		return err
	}

	add := w.stream.Add(*img)

	gocv.PutText(add, fmt.Sprintf("Count: %d 日本語", counter), image.Pt(10, 20), gocv.FontHersheyPlain, 1.2, color.RGBA{0, 255, 0, 0}, 2)
	counter++

	win.IMShow(*add)
	win.WaitKey(1)

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
