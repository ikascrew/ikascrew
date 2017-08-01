package client

import (
	"fmt"
	"time"

	"github.com/ikascrew/ikascrew/video"
)

func (ika *IkascrewClient) xboxController(idx int) error {

	num := idx % 5

	v, err := video.Get(video.Type("file"), fmt.Sprintf("wire/%d.mp4", num))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	fmt.Println("Push " + v.Source())
	err = ika.window.Push(v)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

var xboxFunc func(int) error

func (ika *IkascrewClient) xboxHandleFunc(fn func(int) error) error {
	xboxFunc = fn
	return nil
}

func (ika *IkascrewClient) xboxListen() error {

	ticker := time.NewTicker(time.Second * 10)
	idx := 1
	for {
		select {
		case <-ticker.C:
			idx++
			err := xboxFunc(idx)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
