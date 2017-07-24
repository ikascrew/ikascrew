package client

import (
	"log"

	"github.com/ikascrew/ikascrew"
	"github.com/ikascrew/ikascrew/video"
	"github.com/ikascrew/xbox"
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

	//ika.startHTTP()

	go func() {
		err := display(rep.Project)
		if err != nil {
			log.Fatal(err)
		}
	}()

	xbox.HandleFunc(ika.controller)
	go func() {
		err := xbox.Listen(0)
		if err != nil {
			log.Fatal(err)
		}
	}()

	ika.window.Play(f)

	return nil
}
