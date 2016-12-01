package ikascrew

import (
	"github.com/secondarykey/go-opencv/opencv"
)

type Queue interface {
	Next() *opencv.IplImage
	Size() int
	Wait() int
	Current() int
	Release()
}

type MainVideo struct {
	V1  Queue
	V2  Queue
	img *opencv.IplImage
}

func NewMainVideo(src string) *MainVideo {

	v1, _ := GetVideo(src)

	//dst := opencv.CreateImage(1024, 576, opencv.IPL_DEPTH_8U, 3)

	rtn := &MainVideo{
		V1: v1,
		//img: dst,
	}
	return rtn
}

func (v *MainVideo) Next() *opencv.IplImage {

	img1 := v.V1.Next()

	/*
		img2 := v.V2.Next()
		count := 100
		curCount := v.V1.Current()

		//if curCount <= count {
		alpha := curCount / count
		opencv.AddWeighted(img1, float64(alpha), img2, float64(1.0-alpha), 0.0, v.img)
		return v.img
	*/

	return img1

}

func (v *MainVideo) Size() int {
	return v.V1.Size()
}

func (v *MainVideo) Wait() int {
	return v.V1.Wait()
}

func (v *MainVideo) Current() int {
	return v.V1.Current()
}

func (v *MainVideo) Release() {
	v.V1.Release()
}
