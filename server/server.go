package server

import (
	"fmt"
	"runtime"

	"github.com/google/gops/agent"

	"github.com/ikascrew/ikascrew"
	"github.com/ikascrew/ikascrew/video"

	//pm "github.com/ikascrew/powermate"
	"github.com/ikascrew/xbox"
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

	v, err := video.Get(video.Type(ikascrew.Config.Default.Type), ikascrew.Config.Default.Name)
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

	go func() {
		xbox.HandleFunc(ika.window.EffectXbox)
		xbox.SetDuration(45)
		ika.window.PowerMate = true
		err := xbox.Listen(0)
		if err != nil {
			ika.window.PowerMate = false
			fmt.Printf("Xbox Controller not supported[%v]\n", err)
		}
	}()

	return win.Play(v)
}
