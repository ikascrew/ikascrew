package video

import (
	"fmt"

	"gocv.io/x/gocv"
)

func init() {
}

type Image struct {
	current int
	name    string
	src     *gocv.Mat
}

func NewImage(f string) (*Image, error) {
	img := Image{
		name:    f,
		current: 0,
	}

	wk := gocv.IMRead(f, gocv.IMReadColor)
	if wk.Empty() {
		return nil, fmt.Errorf("Error:LoadImage[%s]", f)
	}

	img.src = &wk

	return &img, nil
}

func (v *Image) Next() (*gocv.Mat, error) {
	v.current++
	if v.current == v.Size() {
		v.current = 0
	}
	return v.src, nil
}

func (v *Image) Set(f int) {
	v.current = f
}

func (v *Image) Current() int {
	return v.current
}

func (v *Image) Size() int {
	return 100
}

func (v *Image) Source() string {
	return v.name
}

func (v *Image) Release() error {
	if !v.src.Empty() {
		v.src.Close()
	}
	return nil
}
