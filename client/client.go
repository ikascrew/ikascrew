package client

import (
	"github.com/ikascrew/ikascrew"
	"github.com/ikascrew/ikascrew/video"
)

func init() {
}

type IkascrewClient struct {
	window *ikascrew.Window
}

func Start() error {

	var err error
	ika := &IkascrewClient{}

	rep, err := ika.syncServer()
	if err != nil {
		return err
	}

	err = ikascrew.Loading(rep.Project)
	if err != nil {
		return err
	}

	f, err := video.Get(video.Type(rep.Type), rep.Source)
	if err != nil {
		return err
	}

	ika.window, err = ikascrew.NewWindow("ikascrew client")
	if err != nil {
		return err
	}

	ika.startHTTP()

	ika.window.Play(f)

	return nil
}
