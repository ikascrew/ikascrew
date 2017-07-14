package server

import (
	"fmt"
	"runtime"

	"github.com/google/gops/agent"

	"github.com/secondarykey/ikascrew"
	"github.com/secondarykey/ikascrew/config"
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

	app, err := config.Get()
	if err != nil {
		return fmt.Errorf("Error:Config[%v]", err)
	}

	v, err := video.Get(video.Type(app.Default.Type), app.Default.Name)
	if err != nil {
		return fmt.Errorf("Error:Video Load[%v]")
	}

	win := ikascrew.NewWindow("ikascrew")
	ika := &IkascrewServer{
		window: win,
	}

	//TODO サーバ起動失敗を見る
	ika.startRPC()

	return win.Play(v)
}
