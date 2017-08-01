package ikascrew

import (
	"fmt"
	"sync"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"

	"github.com/ikascrew/go-opencv/opencv"
	pm "github.com/ikascrew/powermate"
	"github.com/ikascrew/xbox"

	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
)

func init() {
}

type Window struct {
	stream *Stream
	window *opencv.Window

	wait chan Video
	m    *sync.Mutex
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
	rtn.wait = make(chan Video, 100)

	return rtn, nil
}

func (w *Window) Push(v Video) error {
	//return w.stream.Push(v)

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

	js, err := xbox.Open(0)
	if err != nil {
		return fmt.Errorf("Joystick open error.[%v]", err)
	}

	exchange(ProjectName())

	driver.Main(func(s screen.Screen) {
		width := 2048
		height := 192
		opt := &screen.NewWindowOptions{
			Title:  "ikascrew movie viewer",
			Width:  width,
			Height: height,
		}

		shiny, err := s.NewWindow(opt)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer shiny.Release()
		win = shiny

		winSize := image.Point{width, height}
		b, err := s.NewBuffer(winSize)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer b.Release()

		for {

			w.m.Lock()
			fmt.Println("for Lock")
			switch e := shiny.NextEvent().(type) {
			case lifecycle.Event:
				if e.To == lifecycle.StageDead {
					return
				}
			case paint.Event:
				fmt.Println("draw")
				draw(b.RGBA(), current)
				fmt.Println("upload")
				shiny.Upload(image.Point{}, b, b.Bounds())
				fmt.Println("publish")
				shiny.Publish()
			}
			select {
			case v := <-w.wait:
				fmt.Println("stream Push")
				err := w.stream.Push(v)
				fmt.Println("stream End")
				if err != nil {
					fmt.Println("Error")
				}
			default:
				w.stream.PrintVideos("Next() Start")
				img, err := w.stream.Next()
				if err != nil {
				}

				if img != nil {

					fmt.Println("ShowImage")
					w.window.ShowImage(img)

					if err != nil {
					}
					opencv.WaitKey(w.stream.Wait())

					err := Controller(js)
					if err != nil {
						fmt.Println("Controller error")
					}

				} else {
					fmt.Println("Next() Image Nil!!!")
				}
			}
			fmt.Println("for unlock")
			w.m.Unlock()
		}
	})

	return fmt.Errorf("Error : Stream is nil")
}

func (w *Window) Effect(e pm.Event) error {
	return w.stream.Effect(e)
}

func (w *Window) EffectXbox(e xbox.Event) error {
	return w.stream.EffectXbox(e)
}

func (w *Window) Destroy() {
	w.stream.Release()
	w.window.Destroy()
}

func (w *Window) FullScreen() {
	w.window.SetProperty(opencv.CV_WND_PROP_FULLSCREEN, float64(opencv.CV_WINDOW_FULLSCREEN))
}

func (w *Window) Now() {
	//TODO 現状の表示を取得
}
