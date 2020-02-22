package server

import (
	"fmt"
	"runtime"

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
	window *Window
}

func Start(d string) error {

	runtime.GOMAXPROCS(runtime.NumCPU())

	err := ikascrew.Load(d)
	if err != nil {
		return fmt.Errorf("Error Loading directory:%s", err)
	}

	v, err := video.Get(video.Type(ikascrew.Config.Default.Type),
		ikascrew.Config.Default.Name)
	if err != nil {
		return fmt.Errorf("Error:Video Load[%v]", err)
	}

	win, err := NewWindow("ikascrew")
	if err != nil {
		return fmt.Errorf("Error:Create New Window[%v]", err)
	}

	ika := &IkascrewServer{
		window: win,
	}

	go func() {
		ika.startRPC()
	}()

	return win.Play(v)
}
