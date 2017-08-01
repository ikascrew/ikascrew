package ikascrew

import (
	"fmt"

	"github.com/ikascrew/ikascrew/config"
	"github.com/ikascrew/ikascrew/video"
)

var project string
var Config *config.AppConfig

func init() {
	project = ""
	Config = nil
}

func ProjectName() string {
	return project
}

func Loading(name string) error {
	if project == "" {
		project = name
		fmt.Println("Project Name:", project)
	}

	conf, err := config.Load(project)
	if err != nil {
		return fmt.Errorf("Error Config[%v]", err)
	}

	Config = conf

	return nil
}

func Start(p string) error {
	return nil
}

func TestMode(p string, n string) error {

	var err error
	err = Loading(p)
	if err != nil {
		return err
	}

	window, err := NewWindow("ikascrew test")
	if err != nil {
		return err
	}

	v, err := video.Get(project, n)
	if err != nil {
		return err
	}

	//return xbox.Listen(0)
	/*
		go func() {
			err := ika.xboxListen()
			if err != nil {
				fmt.Println("Not Support Xbox")
			}
		}()
	*/

	return window.Play(v)
}
