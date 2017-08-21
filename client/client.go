package client

import (
	"fmt"
	"log"

	"github.com/ikascrew/ikascrew"
	"github.com/ikascrew/xbox"

	"github.com/golang/glog"
	"github.com/google/gops/agent"

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
	err = agent.Listen(nil)
	if err != nil {
		return err
	}

	ika := &IkascrewClient{}
	rep, err := ika.syncServer()
	if err != nil {
		return err
	}

	err = ikascrew.Load(rep.Project)
	if err != nil {
		return err
	}

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

	xbox.HandleFunc(ika.controller)
	err = xbox.Listen(0)
	if err != nil {
		glog.Error("Xbox Listen Error[" + err.Error() + "]")
		return err
	}
	return err
}

var X bool
var A bool
var B bool
var L2 bool
var R2 bool

func (ika *IkascrewClient) controller(e xbox.Event) error {

	if xbox.JudgeAxis(e, xbox.JOY_L_VERTICAL) {
		fmt.Println("JOY L")
		ika.selector.setCursor(e.Axes[xbox.JOY_L_VERTICAL])
		ika.selector.win.Send(paint.Event{})
	}

	if xbox.JudgeAxis(e, xbox.JOY_R_HORIZONTAL) {
		fmt.Println("JOY R")
		ika.pusher.setCursor(e.Axes[xbox.JOY_R_HORIZONTAL])
		ika.pusher.win.Send(paint.Event{})
	}

	if xbox.JudgeAxis(e, xbox.L2) && L2 {
		L2 = false
		err := ika.callPrev()
		if err != nil {
			glog.Error("callPrev[" + err.Error() + "]")
		}
	} else if !xbox.JudgeAxis(e, xbox.L2) {
		L2 = true
	}

	if xbox.JudgeAxis(e, xbox.R2) && R2 {
		R2 = false
		err := ika.callNext()
		if err != nil {
			glog.Error("callNext[" + err.Error() + "]")
		}
	} else if !xbox.JudgeAxis(e, xbox.R2) {
		R2 = true
	}

	if e.Buttons[xbox.X] && X {
		X = false

		res := ika.pusher.get()
		if res != "" {
			err := ika.callEffect(res, "file")
			if err != nil {
				glog.Error("callEffect[" + err.Error() + "]")
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

	if e.Buttons[xbox.B] && B {
		B = false

		err := ika.pusher.delete()
		if err != nil {
			// TODO 無視
			glog.Error("Pusher Delete Error:", err)
		}

	} else if !e.Buttons[xbox.B] {
		B = true
	}

	return nil
}
