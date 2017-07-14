package ikascrew

import (
	"github.com/secondarykey/go-opencv/opencv"
)

type Video interface {
	Next() (*opencv.IplImage, error)
	Wait() int
	Set(int)

	Current() int
	Size() int
	Source() string

	Release() error
}
