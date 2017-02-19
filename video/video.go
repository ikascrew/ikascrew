package video

import (
	"github.com/secondarykey/go-opencv/opencv"
)

type Video interface {
	Next() *opencv.IplImage
	Wait() int
	Size() int
	Current() int
	Set(f int)
	Reload()
	Release()
	Source() string
}
