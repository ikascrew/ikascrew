package client

import (
	"fmt"
	"image"

	"golang.org/x/exp/shiny/screen"
)

type Widget interface {
	Init(screen.Window, screen.Screen, image.Rectangle) error
	Draw()
	Redraw()
	Release()
}

type Part struct {
	owner  screen.Window
	buffer screen.Buffer

	rect image.Rectangle
}

func (p *Part) Init(w screen.Window, s screen.Screen, r image.Rectangle) error {

	bufSize := image.Point{r.Max.X - r.Min.X, r.Max.Y - r.Min.Y}

	b, err := s.NewBuffer(bufSize)
	if err != nil {
		return err
	}
	p.owner = w
	p.buffer = b
	p.rect = r
	return nil
}

func (p *Part) Push() {
	p.owner.Send(p)
}

func (p *Part) Redraw() {

	fmt.Println("Redraw()")

	p.owner.Upload(p.rect.Min, p.buffer, p.buffer.Bounds())
	p.owner.Publish()
}

func (p *Part) Release() {
	p.buffer.Release()
}
