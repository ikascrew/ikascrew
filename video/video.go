package video

import (
	"fmt"
	"github.com/ikascrew/ikascrew"
)

type Type string

const (
	FILE      Type = "file"
	IMAGE     Type = "image"
	MIC       Type = "mic"
	COUNTDOWN Type = "countdown"
	TERMINAL  Type = "terminal"
	LOOP      Type = "loop"
)

//var maps map[string]VideoFactory

func Get(t Type, n string) (ikascrew.Video, error) {

	path := n

	var v ikascrew.Video
	var err error

	switch t {
	case FILE:
		v, err = NewFile(path)
	case COUNTDOWN:
		v, err = NewCountdown(path)
	case TERMINAL:
		v, err = NewTerminal(path)
	case IMAGE:
		v, err = NewImage(path)
	case LOOP:
		v, err = NewLoopFile(path)
	case MIC:
		//v, err = NewMicrophone()
	default:
		err = fmt.Errorf("Not Support Type[%s]", t)
	}

	return v, err
}
