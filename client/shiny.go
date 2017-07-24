package client

import (
	//"github.com/ikascrew/xbox"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

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

	"github.com/ikascrew/ikascrew/video"
	"github.com/ikascrew/xbox"
)

var current image.Image
var idx int
var win screen.Window
var images []string
var resources []string

func init() {
}

func (ika *IkascrewClient) controller(e xbox.Event) error {

	if xbox.JudgeAxis(e, xbox.CROSS_VERTICAL) {
		if e.Axes[xbox.CROSS_VERTICAL] > 0 {
			idx++
		} else {
			idx--
		}
		if idx < 0 {
			idx = 0
		} else if idx >= len(images) {
			idx = len(images) - 1
		}

		file := images[idx]
		src, err := load(file)
		if err != nil {
			return err
		}

		current = src
		win.Send(paint.Event{})
	}

	if e.Buttons[xbox.X] {

		err := ika.callSwitch(resources[idx], "file")
		if err != nil {
			fmt.Println(err)
			return nil
		}
	}

	if e.Buttons[xbox.A] {
		v, err := video.Get(video.Type("file"), resources[idx])
		if err != nil {
			fmt.Println(err)
			return nil
		}

		err = ika.window.Push(v)
		if err != nil {
			fmt.Println(err)
			return nil
		}
	}

	return nil
}

func exchange(dir string) {

	work := dir + "/.public/thumb"

	paths, err := search(work)
	if err != nil {
		panic(err)
	}
	images = paths

	resources = make([]string, len(paths))
	for idx, path := range images {

		jpg := strings.Replace(path, work, "", -1)
		mpg := strings.Replace(jpg, ".jpg", ".mp4", -1)
		path := mpg
		resources[idx] = path
	}

}

func display(dir string) error {

	exchange(dir)

	driver.Main(func(s screen.Screen) {

		width := 2048
		height := 192

		opt := &screen.NewWindowOptions{
			Title:  "ikascrew movie viewer",
			Width:  width,
			Height: height,
		}

		w, err := s.NewWindow(opt)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer w.Release()
		win = w

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
				draw(b.RGBA(), current)
				w.Upload(image.Point{}, b, b.Bounds())
				w.Publish()
			}
		}
	})

	return nil
}

func load(filename string) (image.Image, error) {
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

func draw(m *image.RGBA, img image.Image) {

	if img == nil {
		return
	}

	b := m.Bounds()
	lox := b.Min.X
	loy := b.Min.Y
	hix := b.Max.X
	hiy := b.Max.Y

	for y := loy; y < hiy; y++ {
		for x := lox; x < hix; x++ {
			m.Set(x, y, img.At(x, y))
		}
	}
}

func search(d string) ([]string, error) {

	fileInfos, err := ioutil.ReadDir(d)
	if err != nil {
		return nil, fmt.Errorf("Error:Read Dir[%s]", d)
	}

	rtn := make([]string, 0)

	for _, f := range fileInfos {
		fname := f.Name()
		if f.IsDir() {
			files, err := search(d + "/" + fname)
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
