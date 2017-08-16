package client

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/golang/glog"

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

const SELECTOR_WIDTH = 512
const SELECTOR_HEIGHT = 500

const SELECTOR_IMAGE_W = 512
const SELECTOR_IMAGE_H = 96
const SELECTOR_CUT = 3

type selector struct {
	win     screen.Window
	cursor  int
	current int

	images    []image.Image
	resources []string
}

func NewSelector(dir string) (sel *selector, err error) {

	sel = &selector{}
	err = sel.initialize(dir)
	if err != nil {
		return
	}

	go func() {
		driver.Main(func(s screen.Screen) {

			width := SELECTOR_WIDTH
			height := SELECTOR_HEIGHT

			opt := &screen.NewWindowOptions{
				Title:  "ikascrew selector",
				Width:  width,
				Height: height,
			}

			w, err := s.NewWindow(opt)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer w.Release()
			sel.win = w

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
					sel.draw(b.RGBA())
					w.Upload(image.Point{}, b, b.Bounds())
					w.Publish()
				}
			}
		})
	}()
	return
}

func (s *selector) get() string {
	if s.current < 0 || s.current >= len(s.resources) {
		return ""
	}
	return s.resources[s.current]
}

func (s *selector) setCursor(d int) {
	s.cursor = s.cursor + d
}

func (s *selector) draw(m *image.RGBA) {

	b := m.Bounds()
	lox := b.Min.X
	loy := b.Min.Y
	hix := b.Max.X
	hiy := b.Max.Y

	ver := s.cursor / 1000

	glog.Info("L[%d][%d]\n", s.cursor, ver)
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
			s.current = idx
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

func (s *selector) initialize(p string) error {

	work := p + "/.tmp/thumb"
	paths, err := s.search(work)
	if err != nil {
		return err
	}

	s.images = make([]image.Image, len(paths))
	s.resources = make([]string, len(paths))

	for idx, path := range paths {

		s.images[idx], _ = s.load(path)

		jpg := strings.Replace(path, work, "", -1)
		mpg := strings.Replace(jpg, ".jpg", ".mp4", -1)
		resource := mpg

		s.resources[idx] = resource
	}
	return nil
}

func (s *selector) load(filename string) (image.Image, error) {
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

func (s *selector) search(d string) ([]string, error) {

	fileInfos, err := ioutil.ReadDir(d)
	if err != nil {
		return nil, fmt.Errorf("Error:Read Dir[%s]", d)
	}
	rtn := make([]string, 0)
	for _, f := range fileInfos {
		fname := f.Name()
		if f.IsDir() {
			files, err := s.search(d + "/" + fname)
			if err != nil {
				return nil, err
			}
			rtn = append(rtn, files...)
		} else {
			midx := strings.LastIndex(fname, ".mp4")
			jidx := strings.LastIndex(fname, ".jpg")
			pidx := strings.LastIndex(fname, ".png")
			if midx == len(fname)-4 ||
				jidx == len(fname)-4 ||
				pidx == len(fname)-4 {
				rtn = append(rtn, d+"/"+fname)
			}
		}
	}
	sort.Strings(rtn)
	return rtn, nil
}
