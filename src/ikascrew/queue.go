package ikascrew

import (
	"fmt"
	"github.com/secondarykey/go-opencv/opencv"
)

type Queue struct {
	V1           Video
	V2           Video
	current      Video
	effect       Video
	img          *opencv.IplImage
	stop         bool
	switchCount  int
	switchNumber int
}

func NewQueue(src string) (*Queue, error) {
	v1, err := GetVideo(src)
	if err != nil {
		return nil, err
	}

	dst := opencv.CreateImage(1024, 576, opencv.IPL_DEPTH_8U, 3)
	rtn := &Queue{
		V1:          v1,
		V2:          nil,
		current:     v1,
		effect:      nil,
		switchCount: -1,
		stop:        true,
		img:         dst,
	}
	return rtn, nil
}

func NewSourceQueue(src string, f int) (*Queue, error) {

	v1, err := GetSource(src)
	if err != nil {
		return nil, err
	}

	dst := opencv.CreateImage(1024, 576, opencv.IPL_DEPTH_8U, 3)
	v1.Set(f)
	rtn := &Queue{
		V1:          v1,
		V2:          nil,
		current:     v1,
		effect:      nil,
		switchCount: -1,
		stop:        true,
		img:         dst,
	}
	return rtn, nil
}

func (q *Queue) EffectSwitch(v Video) error {

	if q.V2 == nil {
		return fmt.Errorf("Switch need V2 video")
	}

	q.switchCount = 1
	q.switchNumber = v.Size()
	q.effect = v
	return nil
}

func (q *Queue) Switch(num int) error {
	if q.V2 == nil {
		return fmt.Errorf("Switch need V2 video")
	}

	q.switchCount = 0
	q.switchNumber = num
	return nil
}

func (q *Queue) Next() *opencv.IplImage {

	img := q.current.Next()

	if q.switchNumber == q.switchCount {

		img2 := q.V2.Next()
		q.V1 = q.V2
		q.V2 = q.current

		q.current = q.V1
		q.switchCount = -1
		q.effect = nil

		fmt.Println("Switch done!")
		return img2

	} else if q.switchCount >= 0 {

		img2 := q.V2.Next()
		alpha := float64(q.switchCount) / float64(q.switchNumber)
		if q.effect != nil {
			effect := q.effect.Next()
			opencv.AddWeighted(effect, float64(alpha), img2, float64(1.0-alpha), 0.0, q.img)
			img2 = q.img
		}

		opencv.AddWeighted(img2, float64(alpha), img, float64(1.0-alpha), 0.0, q.img)
		q.switchCount++

		return q.img
	} else if q.effect != nil {
		img2 := q.effect.Next()

		alpha := float64(q.effect.Current()) / float64(q.effect.Size())
		if q.effect.Current() > q.effect.Size()/2 {
			alpha = float64(q.effect.Size()-q.effect.Current()) / float64(q.effect.Size())
		}

		opencv.AddWeighted(img2, float64(alpha), img, float64(1.0-alpha), 0.0, q.img)

		if q.effect.Current() == q.effect.Size() {
			q.effect = nil
		}
		return q.img
	}

	return img

}

func (q *Queue) Effect(v Video) {
	q.effect = v
	return
}

func (q *Queue) Sub(v Video) {
	q.V2 = v
	return
}

func (q *Queue) Name() (string, string) {
	v1Name := ""
	v2Name := "Not Exists"

	for k, v := range videos {
		if v == q.V1 {
			v1Name = k
		} else if q.V2 != nil && v == q.V2 {
			v2Name = k
		}
	}
	return v1Name, v2Name
}

func (q *Queue) Set(v Video, f int) {

	v.Set(f)
	q.V1, q.V2, q.current = v, q.current, v

	return
}

func (q *Queue) Size() int {
	return q.current.Size()
}

func (q *Queue) Wait() int {
	return q.current.Wait()
}

func (q *Queue) Current() int {
	return q.current.Current()
}

func (q *Queue) Release() {
	q.current.Release()
}

func (q *Queue) Source() string {
	return q.current.Source()
}
