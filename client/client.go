package client

import (
	"fmt"
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

	ika.window, err = ikascrew.NewWindow("ikascrew client")
	if err != nil {
		return err
	}
	v, err := video.Get(video.Type(rep.Type), rep.Source)
	if err != nil {
		return err
	}

	go func() {
		err := display(rep.Project)
		if err != nil {
			log.Fatal(err)
		}
	}()

	xbox.HandleFunc(ika.controller)
	//return xbox.Listen(0)
	go func() {
		err := xbox.Listen(0)
		if err != nil {
			log.Fatal(err)
		}
	}()

	return ika.window.Play(v)
}

func TestMode(p string, n string) error {

	var err error
	ika := &IkascrewClient{}

	err = ikascrew.Loading(p)
	if err != nil {
		return err
	}

	ika.window, err = ikascrew.NewWindow("ikascrew client test")
	if err != nil {
		return err
	}

	v, err := video.Get(video.Type("file"), n)
	if err != nil {
		return err
	}

	go func() {
		err := display(p)
		if err != nil {
			log.Fatal(err)
		}
	}()

	xbox.HandleFunc(ika.controller)
	//return xbox.Listen(0)
	go func() {
		err := xbox.Listen(0)
		if err != nil {
			fmt.Println("Not Support Xbox")
		}
	}()

	ika.xboxHandleFunc(ika.xboxController)
	//return xbox.Listen(0)
	/*
		go func() {
			err := ika.xboxListen()
			if err != nil {
				fmt.Println("Not Support Xbox")
			}
		}()
	*/

	return ika.window.Play(v)
}
