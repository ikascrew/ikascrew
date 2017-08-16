package ikascrew

import (
	"github.com/ikascrew/go-opencv/opencv"
)

type Video interface {
	Next() (*opencv.IplImage, error)
	Set(int)

	Current() int
	Size() int
	Source() string

	Release() error
}
