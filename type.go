package ikascrew

import (
	"gocv.io/x/gocv"
)

type Video interface {
	Next() (*gocv.Mat, error)
	Set(int)

	Current() int
	Size() int
	Source() string

	Release() error
}

type Effect interface {
	Run(Video, Video) Video
}
