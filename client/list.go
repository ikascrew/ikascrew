package client

import (
	"image"
	"image/color"
	_ "image/jpeg"

	"strings"

	"golang.org/x/exp/shiny/screen"

	"github.com/ikascrew/ikascrew/tool"
)

var max = 0

type List struct {
	cursor   int
	idx      int
	images   []image.Image
	resource []string
	*Part
}

func NewList(w screen.Window, s screen.Screen) (*List, error) {

	l := &List{}

	r := image.Rect(0, 0, 512, 720)
	l.Part = &Part{}
	l.Init(w, s, r)

	work := "./.client/thumb"

	paths, err := tool.Search(work, nil)
	if err != nil {
		return nil, err
	}

	l.images = make([]image.Image, len(paths)+1)
	l.resource = make([]string, len(paths)+1)

	for idx, path := range paths {
		l.images[idx], _ = tool.LoadImage(path)
		jpg := strings.Replace(path, work, "", -1)
		mpg := strings.Replace(jpg, ".jpg", ".mp4", -1)
		resource := mpg
		l.resource[idx] = resource
	}

	max = len(paths) * 100 * 100

	return l, nil
}

func (l *List) Draw() {

	m := l.Part.buffer.RGBA()

	lox := 0
	loy := 0
	hix := 512
	hiy := 720

	ver := l.cursor / 200

	white := color.RGBA{255, 255, 255, 255}
	black := color.RGBA{0, 0, 0, 255}

	start := (ver / 96)

	for y := loy; y < hiy; y++ {

		var img image.Image

		d := y / 96
		idx := start + d

		if idx >= 0 && idx < len(l.images) {
			img = l.images[start+d]
		}

		dy := y - (d * 96)

		flag := false
		yflag := false

		if (y+48) > 240 && (y-48) < 240 {
			l.idx = idx
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

func (l *List) get() string {
	if l.idx < 0 || l.idx >= len(l.resource) {
		return ""
	}
	return l.resource[l.idx]
}

func (l *List) setCursor(d int) {
	l.cursor = l.cursor + d
	l.Draw()
}

func (l *List) zeroCursor() {
	l.cursor = 0
	l.Draw()
}

func (l *List) maxCursor() {
	l.cursor = l.cursor + max
	l.Draw()
}
