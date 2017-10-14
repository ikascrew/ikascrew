package client

import (
	"github.com/ikascrew/ikascrew"
	"github.com/ikascrew/ikascrew/tool"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"

	"github.com/golang/glog"

	"golang.org/x/mobile/event/paint"

	"fmt"
	"image/color"
	"os"
	"strings"

	"image"
	_ "image/jpeg"
)

type Window struct {
	window screen.Window
	buffer screen.Buffer

	cursorS   int
	currentS  int
	images    []image.Image
	resourceS []string

	cursorP   int
	currentP  int
	targets   []image.Image
	resourceP []string
}

func NewWindow(t string, w, h int) (window *Window, err error) {

	window = &Window{
		window: nil,
		buffer: nil,
	}

	driver.Main(func(s screen.Screen) {
		opt := &screen.NewWindowOptions{
			Title:  t,
			Width:  w,
			Height: h,
		}
		window.window, err = s.NewWindow(opt)
		if err != nil {
			return
		}
		winSize := image.Point{w, h}
		window.buffer, err = s.NewBuffer(winSize)
		if err != nil {
			return
		}
	})

	err = window.initialize(ikascrew.ProjectName())
	return
}

func (w *Window) Repaint() {
	w.window.Send(paint.Event{})
}

func (w *Window) Release() {
	w.window.Release()
	w.buffer.Release()
}

func (s *Window) drawSelector() {

	m := s.buffer.RGBA()

	lox := 0
	loy := 0
	hix := 512
	hiy := 720

	ver := s.cursorS / 200

	glog.Info("L[%d][%d]\n", s.cursorS, ver)
	white := color.RGBA{255, 255, 255, 255}
	black := color.RGBA{0, 0, 0, 255}

	start := (ver / 96)

	for y := loy; y < hiy; y++ {

		var img image.Image

		d := y / 96
		idx := start + d

		if idx >= 0 && idx < len(s.images) {
			img = s.images[start+d]
		}

		dy := y - (d * 96)

		flag := false
		yflag := false

		if (y+48) > 240 && (y-48) < 240 {
			s.currentS = idx
			if dy <= 5 || dy >= 91 {
				flag = true
			} else {
				yflag = true
			}
		}

		for x := lox; x < hix; x++ {

			if yflag {
				if x <= 5 || x >= 507 {
					flag = true
				} else {
					flag = false
				}
			}

			go func(img image.Image, x, y, dy int, flag bool) {
				if img == nil {
					m.Set(x, y, black)
				} else if flag {
					m.Set(x, y, white)
				} else {
					m.Set(x, y, img.At(x, dy))
				}
			}(img, x, y, dy, flag)
		}
	}
}

func (s *Window) getS() string {
	if s.currentS < 0 || s.currentS >= len(s.resourceS) {
		return ""
	}
	return s.resourceS[s.currentS]
}

func (s *Window) setCursorS(d int) {
	s.cursorS = s.cursorS + d
	s.drawSelector()
}

func (s *Window) initialize(p string) error {

	work := p + "/.client/thumb"
	paths, err := tool.Search(work, nil)
	if err != nil {
		return err
	}

	s.images = make([]image.Image, len(paths))
	s.resourceS = make([]string, len(paths))

	for idx, path := range paths {

		s.images[idx], _ = s.load(path)

		jpg := strings.Replace(path, work, "", -1)
		mpg := strings.Replace(jpg, ".jpg", ".mp4", -1)
		resource := mpg

		s.resourceS[idx] = resource
	}
	return nil
}

func (s *Window) load(filename string) (image.Image, error) {
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

func (p *Window) getP() string {

	sz := len(p.resourceP)
	if p.currentP > sz-1 {
		return ""
	}

	rtn := p.resourceP[p.currentP]

	err := p.delete()
	if err != nil {
		glog.Error(err.Error())
		return ""
	}

	return rtn
}

func (p *Window) delete() error {

	sz := len(p.resourceP)
	if p.currentP > sz-1 {
		return fmt.Errorf("Pusher Index Error")
	}

	newres := make([]string, 0)
	newtar := make([]image.Image, 0)
	for idx, elm := range p.resourceP {
		if idx != p.currentP {
			newres = append(newres, elm)
			newtar = append(newtar, p.targets[idx])
		}
	}
	p.resourceP = newres
	p.targets = newtar

	p.cursorP = 0

	p.drawPusher()
	p.Repaint()
	return nil
}

func (p *Window) add(f string) error {

	for _, elm := range p.resourceP {
		if f == elm {
			return fmt.Errorf("Resource[" + f + "] exist")
		}
	}
	p.resourceP = append(p.resourceP, f)

	icon := strings.Replace(f, ".mp4", ".jpg", 1)

	d := ikascrew.ProjectName()
	file := d + "/.client/icon" + icon
	img, err := p.load(file)
	if err != nil {
		return err
	}
	p.targets = append(p.targets, img)

	p.cursorP = 0

	p.drawPusher()
	p.Repaint()
	return nil
}

func (p *Window) drawPusher() {

	m := p.buffer.RGBA()

	lox := 512
	loy := 0
	hix := 1536
	hiy := 144

	hor := p.cursorP / 175
	glog.Info("R[%d][%d]\n", p.cursorP, hor)

	white := color.RGBA{255, 255, 255, 255}
	black := color.RGBA{0, 0, 0, 255}

	start := (hor / 256)

	for y := loy; y < hiy; y++ {
		var img image.Image
		for x := lox; x < hix; x++ {

			d := x / 256
			idx := start + d

			flag := false
			if x >= 512 && x < 768 {
				p.currentP = idx
				flag = true
				if x > 517 && x < 763 {
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

func (s *Window) setCursorP(d int) {
	s.cursorP = s.cursorP + d
	s.drawPusher()
}

func (w *Window) keyListener(k int) {

	fmt.Printf("[%d]\n", k)

	switch k {
	case 82:
		w.setCursorS(-60000)
		w.Repaint()
	case 81:
		w.setCursorS(60000)
		w.Repaint()
	case 80:
		w.setCursorP(10000)
	case 79:
		w.setCursorP(-10000)
	case 40:
	}
}
