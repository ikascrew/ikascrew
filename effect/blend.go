package effect

import (
	"github.com/secondarykey/go-opencv/opencv"
	"github.com/secondarykey/ikascrew"
)

type Blend struct {
	s *Switch
}

func NewBlend(v ikascrew.Video, e ikascrew.Effect) (*Blend, error) {

	s, err := NewSwitch(v, e)
	if err != nil {
		return nil, err
	}
	b := Blend{
		s: s,
	}
	return &b, nil
}

func (b *Blend) Next() (*opencv.IplImage, error) {

	now, _ := b.s.now.Next()
	next, _ := b.s.video.Next()

	alpha := b.s.count / b.s.number

	opencv.AddWeighted(next, float64(alpha), now, float64(1.0-alpha), 0.0, b.s.img)

	if alpha <= 0.5 {
		b.s.count++
	}
	return b.s.img, nil
}

func (b *Blend) Wait() int {
	return b.s.Wait()
}

func (b *Blend) Release() error {
	return b.s.Release()
}

func (b *Blend) String() string {
	return b.s.String()
}
