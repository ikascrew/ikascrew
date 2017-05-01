package ikascrew

import (
	"github.com/secondarykey/go-opencv/opencv"
)

type Video interface {
	//Load(string) error

	Next() (*opencv.IplImage, error)
	Wait() int
	Set(int)

	Current() int
	Size() int
	Source() string

	Release() error
}

type Effect interface {

	//	Change(Effect) error
	Next() (*opencv.IplImage, error)

	Wait() int
	Release() error

	String() string
}
