package server

import (
	"fmt"
	"runtime"

	"github.com/google/gops/agent"

	"github.com/secondarykey/ikascrew"
	"github.com/secondarykey/ikascrew/video"

	pm "github.com/secondarykey/powermate"
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

	//TODO サーバ起動失敗を見る
	ika.startRPC()

	//Effect Handling
	go func() {
		pm.HandleFunc(ika.window.Effect)
		ika.window.PowerMate = true
		err := pm.Listen("/dev/input/powermate")
		if err != nil {
			ika.window.PowerMate = false
			fmt.Printf("Powermate not supported[%v]\n", err)
		}
	}()

	return win.Play(v)
}
