package server

import (
	"fmt"
	"runtime"

	"github.com/google/gops/agent"

	"github.com/secondarykey/ikascrew"
	"github.com/secondarykey/ikascrew/effect"
	"github.com/secondarykey/ikascrew/video"
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

	if err := agent.Listen(nil); err != nil {
		return err
	}

	runtime.GOMAXPROCS(runtime.NumCPU())
	err := ikascrew.Loading(d)

	if err != nil {
		return fmt.Errorf("Error Loading directory:%s", err)
	}

	//TODO app.json をみよう
	f, err := video.NewImage(d + "/logo.png")
	if err != nil {
		return fmt.Errorf("Error:Video Load")
	}

	e, err := effect.NewNormal(f)
	if err != nil {
		return fmt.Errorf("Error:Effect")
	}

	win := ikascrew.NewWindow("ikascrew", e)
	ika := &IkascrewServer{
		window: win,
	}

	ika.startRPC()

	return win.Play()
}
