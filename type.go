package ikascrew

import (
	"gocv.io/x/gocv"
	//"github.com/ikascrew/go-opencv/opencv"
)

type Video interface {
	Next() (*gocv.Mat, error)
	Set(int)

	Current() int
	Size() int
	Source() string

	Release() error
}
