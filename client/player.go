package client

import (
	//"github.com/ikascrew/xbox"
	"fmt"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"

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

type player struct {
	win    screen.Window
	target string
}

func NewPlayer() (p *player, err error) {

	p = &player{}

	go func() {
		driver.Main(func(s screen.Screen) {

			width := 1024
			height := 576

			opt := &screen.NewWindowOptions{
				Title:  "ikascrew player",
				Width:  width,
				Height: height,
			}

			w, err := s.NewWindow(opt)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer w.Release()
			p.win = w

			winSize := image.Point{width, height}
			b, err := s.NewBuffer(winSize)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer b.Release()

			for {
				switch e := w.NextEvent().(type) {
				case lifecycle.Event:
					if e.To == lifecycle.StageDead {
						return
					}
				case paint.Event:
					//draw(b.RGBA(), current)
					//w.Upload(image.Point{}, b, b.Bounds())
					w.Publish()
				}
			}
		})
	}()
	return
}
