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

	err := ika.syncServer()
	if err != nil {
		return err
	}

	err = ikascrew.Loading("projects/20170502")
	if err != nil {
		return err
	}

	f, err := video.NewImage("projects/20170502/logo.png")

	n, err := effect.NewNormal(f)
	if err != nil {
		return err
	}

	ika.window = ikascrew.NewWindow("ikascrew client", n)

	ika.startHTTP()
	ika.window.Play()
	return nil
}
