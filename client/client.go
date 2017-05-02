package client

import (
	"github.com/secondarykey/ikascrew"
	"github.com/secondarykey/ikascrew/effect"
	"github.com/secondarykey/ikascrew/video"
)

func init() {

}

type IkascrewClient struct {
	window *ikascrew.Window
}

func Start() error {

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

	n, err := effect.NewNormal(f)
	if err != nil {
		return err
	}

	ika.window = ikascrew.NewWindow("ikascrew client", n)

	ika.startHTTP()
	ika.window.Play()
	return nil
}
