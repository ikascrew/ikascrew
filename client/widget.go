package client

import (
	"image"

	"golang.org/x/exp/shiny/screen"
)

type Widget interface {
	Init(screen.Window, screen.Screen, image.Rectangle) error
	Draw()
	Redraw()
}

type Part struct {
	owner  screen.Window
	buffer screen.Buffer

	rect image.Rectangle
}

func (p *Part) Init(w screen.Window, s screen.Screen, r image.Rectangle) error {

	b, err := s.NewBuffer(r.Max)
	if err != nil {
		return err
	}
	p.owner = w
	p.buffer = b
	p.rect = r
	return nil
	//p.owner.Send(paint.Event{})
}

func (p *Part) Redraw() {
	p.owner.Upload(p.rect.Min, p.buffer, p.buffer.Bounds())
	p.owner.Publish()
}
