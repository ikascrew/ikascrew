package server

import (
	"fmt"
	"runtime"

	"github.com/golang/glog"

	"github.com/google/gops/agent"

	"github.com/ikascrew/ikascrew"
	"github.com/ikascrew/ikascrew/video"

	pm "github.com/ikascrew/powermate"
)

func init() {
}

const ADDRESS = ":55555"

func Address() string {
	return ADDRESS
}

type IkascrewServer struct {
	window *ikascrew.Window
}

func Start(d string) error {

	glog.Info("Start gops Agent")
	var err error
	err = agent.Listen(nil)
	if err != nil {
		return err
	}

	glog.Info("Set max procs")
	runtime.GOMAXPROCS(runtime.NumCPU())

	err = ikascrew.Loading(d)
	if err != nil {
		return fmt.Errorf("Error Loading directory:%s", err)
	}

	glog.Info("Set Default video")
	v, err := video.Get(video.Type(ikascrew.Config.Default.Type),
		ikascrew.Config.Default.Name)
	if err != nil {
		return fmt.Errorf("Error:Video Load[%v]", err)
	}

	glog.Info("Create main window")
	win, err := ikascrew.NewWindow("ikascrew")
	if err != nil {
		return fmt.Errorf("Error:Create New Window[%v]", err)
	}

	ika := &IkascrewServer{
		window: win,
	}

	glog.Info("Initialize powermate")
	go func() {
		pm.HandleFunc(ika.window.Effect)
		ika.window.PowerMate = true
		err := pm.Listen("/dev/input/powermate")
		if err != nil {
			ika.window.PowerMate = false
			glog.Error("Powermate not supported[", err, "]")
		}
	}()

	glog.Info("Start RPC")
	ika.startRPC()

	glog.Info("Let's Play!")
	return win.Play(v)
}
