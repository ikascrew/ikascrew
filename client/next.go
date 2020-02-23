package client

import (
	"fmt"
	"image"
	"image/color"
	"strings"

	"golang.org/x/exp/shiny/screen"

	"github.com/ikascrew/ikascrew/tool"
)

type Next struct {
	cursor   int
	idx      int
	targets  []image.Image
	resource []string

	*Part
}

func NewNext(w screen.Window, s screen.Screen) (*Next, error) {

	n := &Next{}
	r := image.Rect(512, 0, 1536, 144)

	n.Part = &Part{}
	n.Init(w, s, r)

	n.targets = make([]image.Image, 0)
	n.resource = make([]string, 0)

	return n, nil
}

func (n *Next) Draw() {

	m := n.Part.buffer.RGBA()

	lox := 0
	loy := 0
	hix := 1024
	hiy := 144

	hor := n.cursor / 175

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
				n.idx = idx
				flag = true
				if x > 5 && x < 251 {
					if y > 5 && y < 140 {
						flag = false
					}
				}
			}

			if idx >= 0 && idx < len(n.targets) {
				img = n.targets[idx]
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
}

func (n *Next) get() string {
	sz := len(n.resource)
	if n.idx > sz-1 || sz == 0 || n.idx < 0 {
		return ""
	}

	rtn := n.resource[n.idx]
	return rtn
}

func (n *Next) delete() error {

	sz := len(n.resource)
	if n.idx > sz-1 || sz == 0 || n.idx < 0 {
		return fmt.Errorf("Pusher Index Error")
	}

	newres := make([]string, 0)
	newtar := make([]image.Image, 0)
	for idx, elm := range n.resource {
		if idx != n.idx {
			newres = append(newres, elm)
			newtar = append(newtar, n.targets[idx])
		}
	}
	n.resource = newres
	n.targets = newtar

	n.cursor = 0
	n.Draw()

	return nil
}

func (n *Next) add(f string) error {

	for _, elm := range n.resource {
		if f == elm {
			return fmt.Errorf("Resource[" + f + "] exist")
		}
	}
	n.resource = append(n.resource, f)

	icon := strings.Replace(f, ".mp4", ".jpg", 1)

	file := "./.client/icon" + icon

	img, err := tool.LoadImage(file)

	if err != nil {
		return err
	}
	n.targets = append(n.targets, img)

	n.cursor = 0
	n.Draw()

	return nil
}

func (n *Next) setCursor(d int) {
	n.cursor = n.cursor + d
	n.Draw()
}
