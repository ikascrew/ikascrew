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

type BlendVideo struct {
	V1  Queue
	V2  Queue
	img *opencv.IplImage
}

func NewBlendVideo(src, dest string) *BlendVideo {

	v1, _ := GetVideo(src)
	v2, _ := GetVideo(dest)
	dst := opencv.CreateImage(1024, 576, opencv.IPL_DEPTH_8U, 3)

	rtn := &BlendVideo{
		V1:  v1,
		V2:  v2,
		img: dst,
	}
	return rtn
}

func (v *BlendVideo) Next() *opencv.IplImage {

	img1 := v.V1.Next()
	img2 := v.V2.Next()

	count := 100
	curCount := v.V1.Current()

	//if curCount <= count {
	alpha := curCount / count
	opencv.AddWeighted(img1, float64(alpha), img2, float64(1.0-alpha), 0.0, v.img)
	return v.img
	//}
	return img2

}

func (v *BlendVideo) Size() int {
	return 0
}

func (v *BlendVideo) Wait() int {
	return v.V1.Wait()
}

func (v *BlendVideo) Current() int {
	return 0
}

func (v *BlendVideo) Release() {
	v.V1.Release()
	v.V2.Release()
}
