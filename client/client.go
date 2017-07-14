package client

import (
	"github.com/secondarykey/ikascrew"
	"github.com/secondarykey/ikascrew/video"
	pm "github.com/secondarykey/powermate"
)

func init() {
}

type IkascrewClient struct {
	window *ikascrew.Window
}

func Start() error {

	ika := &IkascrewClient{}

	rep, err := ika.syncServer()
	if err != nil {
		return err
	}

	err = ikascrew.Loading(rep.Project)
	if err != nil {
		return err
	}

	f, err := video.Get(video.Type(rep.Type), rep.Source)
	if err != nil {
		return err
	}

	ika.window = ikascrew.NewWindow("ikascrew client")

	ika.startHTTP()

	//Effect Handling
	go func() {
		pm.HandleFunc(ika.window.Effect)
		err := pm.Listen("/dev/input/powermate")
		if err != nil {
			fmt.Println("Powermate not supported")
		} else {
			ika.window.PowerMate = true
		}
	}()

	ika.window.Play(f)

	return nil
}
