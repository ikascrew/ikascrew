package server

import (
	"fmt"
	"runtime"

	"github.com/google/gops/agent"

	"github.com/ikascrew/ikascrew"
	"github.com/ikascrew/ikascrew/video"
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

	var err error
	err = agent.Listen(nil)
	if err != nil {
		return err
	}

	runtime.GOMAXPROCS(runtime.NumCPU())
	err = ikascrew.Loading(d)

	if err != nil {
		return fmt.Errorf("Error Loading directory:%s", err)
	}

	v, err := video.Get(d, ikascrew.Config.Default.Name)
	if err != nil {
		return fmt.Errorf("Error:Video Load[%v]", err)
	}

	win, err := ikascrew.NewWindow("ikascrew")
	if err != nil {
		return fmt.Errorf("Error:Create New Window[%v]", err)
	}

	ika := &IkascrewServer{
		window: win,
	}

	ika.startRPC()

	return win.Play(v)
}
