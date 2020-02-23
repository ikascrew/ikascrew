package client

import (
	"os"
	"strconv"

	"github.com/ikascrew/ikascrew"
	"github.com/ikascrew/ikascrew/pb"
	pm "github.com/ikascrew/powermate"
	vol "github.com/ikascrew/volumes"
	"github.com/ikascrew/xbox"

	"github.com/golang/glog"
)

var vols *vol.Volumes

func init() {
	vols = vol.New()
	vols.Add("Volume", 300)
	vols.Add("Light", 100)
	vols.Add("Wait", 50)
}

type IkascrewClient struct {
	selector *Window
	testMode bool
}

var ika *IkascrewClient

func Start() error {

	var err error
	var rep *pb.SyncReply

	go vols.Start()
	pm.HandleFunc(trigger)

	ika = &IkascrewClient{}

	args := os.Args
	if len(args) > 2 {
		ika.testMode = true
	} else {
		ika.testMode = false
	}

	if ika.testMode {
		pid, _ := strconv.Atoi(args[2])
		rep = &pb.SyncReply{
			Source:  "",
			Type:    "",
			Project: int64(pid),
		}
	} else {
		rep, err = ika.syncServer()
		if err != nil {
			return err
		}
	}

	err = ikascrew.Load(rep.Project)
	if err != nil {
		return err
	}

	//XBOX Controller
	xbox.HandleFunc(ika.controller)
	go func() {
		err = xbox.Listen(0)
		if err != nil {
			glog.Error("Xbox Listen Error[" + err.Error() + "]")
			return
		}
	}()

	//powermate
	go func() {
		err = pm.Listen("/dev/input/powermate")
		if err != nil {
			glog.Error("powermate Listen Error[" + err.Error() + "]")
			return
		}
	}()

	//Main
	win, err := NewWindow("ikascrew client", 1536, 766)
	if err != nil {
		glog.Error("NewWindow() Error[" + err.Error() + "]")
		return err
	}

	ika.selector = win
	win.SetClient(ika)

	//クライアント描画
	for {
		e := win.window.NextEvent()
		switch e := e.(type) {
		case *Part:
			e.Redraw()
		}
	}
	return err
}

func trigger(e pm.Event) error {

	val := vols.Get()
	if zero {
		val = 0
		vols.SetCursor(0)
		zero = false
	}
	idx := vols.GetCursor()
	update := false

	switch e.Type {
	case pm.Rotation:
		switch e.Value {
		case pm.Left:
			val -= 1.0
		case pm.Right:
			val += 1.0
		}

		update = true
	case pm.Press:
		switch e.Value {
		case pm.Up:
		case pm.Down:
			idx = idx + 1
			if idx > 2 {
				idx = 0
			}
			vols.SetCursor(idx)
		}
	default:
	}

	if update {

		vols.Set(val)
		var i int64 = int64(idx)

		message := pb.VolumeMessage{
			Index: i,
			Value: val,
		}

		err := ika.callVolume(message)
		if err != nil {
			return err
		}

		//Reply で再度設定
	}

	return nil
}

var zero = false

func setZero() {
	zero = true
}
