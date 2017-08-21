package client

import (
	"fmt"
	"os"
	"strings"

	"github.com/golang/glog"
	"github.com/ikascrew/ikascrew"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"

	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

type pusher struct {
	win     screen.Window
	cursor  int
	current int

	targets   []image.Image
	resources []string
}

func NewPusher() (p *pusher, err error) {

	p = &pusher{
		targets:   make([]image.Image, 0),
		resources: make([]string, 0),
	}

	go func() {
		driver.Main(func(s screen.Screen) {

			width := 1024
			height := 144

			opt := &screen.NewWindowOptions{
				Title:  "ikascrew pusher",
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
					p.draw(b.RGBA())
					w.Upload(image.Point{}, b, b.Bounds())
					w.Publish()
				}
			}
		})
	}()
	return
}

func (p *pusher) get() string {

	sz := len(p.resources)
	if p.current > sz-1 {
		return ""
	}

	rtn := p.resources[p.current]

	err := p.delete()
	if err != nil {
		glog.Error(err.Error())
		return ""
	}

	return rtn
}

func (p *pusher) delete() error {

	sz := len(p.resources)
	if p.current > sz-1 {
		return fmt.Errorf("Pusher Index Error")
	}

	newres := make([]string, 0)
	newtar := make([]image.Image, 0)
	for idx, elm := range p.resources {
		if idx != p.current {
			newres = append(newres, elm)
			newtar = append(newtar, p.targets[idx])
		}
	}
	p.resources = newres
	p.targets = newtar
	p.win.Send(paint.Event{})

	p.cursor = 0
	p.win.Send(paint.Event{})
	return nil
}

func (p *pusher) add(f string) error {

	for _, elm := range p.resources {
		if f == elm {
			return fmt.Errorf("Resource[" + f + "] exist")
		}
	}
	p.resources = append(p.resources, f)

	icon := strings.Replace(f, ".mp4", ".jpg", 1)

	d := ikascrew.ProjectName()
	file := d + "/.client/icon" + icon
	img, err := p.load(file)
	if err != nil {
		return err
	}
	p.targets = append(p.targets, img)

	p.cursor = 0
	p.win.Send(paint.Event{})
	return nil
}

func (s *pusher) load(filename string) (image.Image, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	m, _, err := image.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("could not decode %s: %v", filename, err)
	}
	return m, nil
}

func (p *pusher) draw(m *image.RGBA) {

	b := m.Bounds()
	lox := b.Min.X
	loy := b.Min.Y
	hix := b.Max.X
	hiy := b.Max.Y

	hor := p.cursor / 175
	glog.Info("R[%d][%d]\n", p.cursor, hor)

	white := color.RGBA{255, 255, 255, 255}
	black := color.RGBA{0, 0, 0, 255}

	start := (hor / 256)

	for y := loy; y < hiy; y++ {
		var img image.Image
		for x := lox; x < hix; x++ {

			d := x / 256
			idx := start + d

			flag := false
			if x >= 0 && x < 256 {
				p.current = idx
				flag = true
				if x > 5 && x < 251 {
					if y > 5 && y < 140 {
						flag = false
					}
				}
			}

			if idx >= 0 && idx < len(p.targets) {
				img = p.targets[idx]
			} else {
				img = nil
			}

			dx := x - (d * 256)
			go func(img image.Image, x, y, dx int, flag bool) {
				if img == nil {
					m.Set(x, y, black)
				} else if flag {
					m.Set(x, y, white)
				} else {
					m.Set(x, y, img.At(dx, y))
				}
			}(img, x, y, dx, flag)
		}
	}
	return
}

func (s *pusher) setCursor(d int) {
	s.cursor = s.cursor + d
}
