package client

import (
	"fmt"
	"log"

	"github.com/ikascrew/ikascrew"
	"github.com/ikascrew/xbox"

	"golang.org/x/mobile/event/paint"
)

func init() {
}

type IkascrewClient struct {
	selector *selector
	player   *player
	pusher   *pusher
}

func Start() error {

	var err error

	//sync
	d := "projects/20170817"

	ika := &IkascrewClient{}
	err = ikascrew.Loading(d)
	if err != nil {
		return err
	}

	selector, err := NewSelector(d)
	if err != nil {
		log.Fatal(err)
	}

	pusher, err := NewPusher()
	if err != nil {
		log.Fatal(err)
	}

	player, err := NewPlayer()
	if err != nil {
		log.Fatal(err)
	}

	ika.selector = selector
	ika.player = player
	ika.pusher = pusher

	xbox.HandleFunc(ika.controller)
	err = xbox.Listen(0)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

var X bool
var A bool

func (ika *IkascrewClient) controller(e xbox.Event) error {

	if xbox.JudgeAxis(e, xbox.JOY_L_VERTICAL) {
		ika.selector.setCursor(e.Axes[xbox.JOY_L_VERTICAL])
		ika.selector.win.Send(paint.Event{})
	}

	if xbox.JudgeAxis(e, xbox.JOY_R_HORIZONTAL) {
		ika.pusher.setCursor(e.Axes[xbox.JOY_R_HORIZONTAL])
		ika.pusher.win.Send(paint.Event{})
	}

	if e.Buttons[xbox.X] && X {
		X = false
		err := ika.callSwitch(ika.pusher.get(), "file")
		if err != nil {
			fmt.Println(err)
		}
	} else if !e.Buttons[xbox.X] {
		X = true
	}

	if e.Buttons[xbox.A] && A {

		A = false
		fmt.Printf("[%s]\n", ika.selector.get())
		err := ika.pusher.add(ika.selector.get())
		if err != nil {
			fmt.Println(err)
		}

	} else if !e.Buttons[xbox.A] {
		A = true
	}

	return nil
}
