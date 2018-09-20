package client

import (
	"github.com/ikascrew/ikascrew"
	"github.com/ikascrew/xbox"

	"github.com/golang/glog"
	"github.com/google/gops/agent"

	"golang.org/x/mobile/event/key"
)

func init() {
}

type IkascrewClient struct {
	selector *Window
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

	xbox.HandleFunc(ika.controller)
	go func() {
		err = xbox.Listen(0)
		if err != nil {
			glog.Error("Xbox Listen Error[" + err.Error() + "]")
			return
		}
	}()

	win, err := NewWindow("ikascrew client", 1536, 720)
	if err != nil {
		glog.Error("NewWindow() Error[" + err.Error() + "]")
		return err
	}

	ika.selector = win

	for {
		e := win.window.NextEvent()
		switch e := e.(type) {
		case key.Event:
			win.keyListener(int(e.Code))
		case *Part:
			e.Redraw()
		}
	}
	return err
}
