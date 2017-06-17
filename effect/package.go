package effect

import (
	"fmt"

	"github.com/secondarykey/ikascrew"
)

type Type string

const (
	SWITCH Type = "switch"
	MATE   Type = "mate"
)

func Get(t Type, v ikascrew.Video, now ikascrew.Effect) (ikascrew.Effect, error) {

	var e ikascrew.Effect
	var err error

	switch t {
	case SWITCH:
		e, err = NewSwitch(v, now)
	case MATE:
		switch now.(type) {
		case *Mate:
			err = fmt.Errorf("MateEffect Can not be used continuously.")
		default:
			e, err = NewMate(v, now)
		}
	default:
		err = fmt.Errorf("Not Support Effect[%s]", t)
	}
	return e, err
}
