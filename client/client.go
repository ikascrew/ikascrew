package client

import (
	"log"

	"github.com/ikascrew/ikascrew"
	"github.com/ikascrew/xbox"

	"golang.org/x/mobile/event/paint"

	"github.com/golang/glog"
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

	ika := &IkascrewClient{}
	rep, err := ika.syncServer()
	if err != nil {
		return err
	}

	err = ikascrew.Load(rep.Project)
	if err != nil {
		return err
	}

	xbox.HandleFunc(ika.controller)
	err = xbox.Listen(0)
	if err != nil {
		glog.Error("Xbox Listen Error[" + err.Error() + "]")
		ika.startHTTP()
	} else {
		selector, err := NewSelector(rep.Project)
		if err != nil {
			log.Fatal(err)
		}
		ika.selector = selector

		pusher, err := NewPusher()
		if err != nil {
			log.Fatal(err)
		}
		ika.pusher = pusher

		/*
			player, err := NewPlayer()
			if err != nil {
				log.Fatal(err)
			}
			ika.player = player
		*/
	}
	return err
}

var X bool
var A bool

func (ika *IkascrewClient) controller(e xbox.Event) error {

	if xbox.JudgeAxis(e, xbox.JOY_L_VERTICAL) {
		ika.selector.setCursor(e.Axes[xbox.JOY_L_VERTICAL])
		ika.selector.win.Send(paint.Event{})
	}

	if xbox.JudgeAxis(e, xbox.JOY_L_HORIZONTAL) {
		ika.pusher.setCursor(e.Axes[xbox.JOY_L_HORIZONTAL])
		ika.pusher.win.Send(paint.Event{})
	}

	if e.Buttons[xbox.X] && X {
		X = false

		res := ika.pusher.get()
		if res != "" {
			err := ika.callSwitch(res, "file")
			if err != nil {
				glog.Error("callSwitch[" + err.Error() + "]")
			}
		} else {
			glog.Error("Pusher Error: No Index")
		}
	} else if !e.Buttons[xbox.X] {
		X = true
	}

	if e.Buttons[xbox.A] && A {
		A = false

		res := ika.selector.get()

		if res != "" {
			err := ika.pusher.add(res)
			if err != nil {
				// TODO 無視
				glog.Error("Pusher Add Error:", err)
			}
		} else {
			glog.Error("Selector Error:" + "No Index")
		}

	} else if !e.Buttons[xbox.A] {
		A = true
	}

	return nil
}
