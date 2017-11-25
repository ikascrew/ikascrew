package video

import (
	"fmt"
	"github.com/ikascrew/ikascrew"
)

type Type string

const (
	FILE  Type = "file"
	IMAGE Type = "image"
	MIC   Type = "mic"
)

func Get(t Type, n string) (ikascrew.Video, error) {

	path := ikascrew.ProjectName() + "/" + n

	var v ikascrew.Video
	var err error

	switch t {
	case FILE:
		v, err = NewFile(path)
	case IMAGE:
		v, err = NewImage(path)
	case MIC:
		//v, err = NewMicrophone()
	default:
		err = fmt.Errorf("Not Support Type[%s]", t)
	}

	return v, err
}
