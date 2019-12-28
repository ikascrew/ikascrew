package server

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/golang/glog"

	"github.com/ikascrew/ikascrew"

	"gocv.io/x/gocv"
)

func init() {
}

type Window struct {
	name string
	wait chan ikascrew.Video

	stream *Stream
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

	//イメージを取得
	img, err := w.stream.Get()
	if err != nil {
		return err
	}
	//作成
	add := w.stream.Add(*img)
	//表示
	win.IMShow(*add)
	win.WaitKey(1)

	wg.Wait()

	return nil
}

func (w *Window) SetSwitch(t string) error {
	return w.stream.SetSwitch(t)
}

func (w *Window) Destroy() {
	w.stream.Release()
}

func (w *Window) FullScreen() {
	//TODO gocvでのフル
}
