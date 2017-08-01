package ikascrew

import (
	"github.com/ikascrew/xbox"
	"golang.org/x/mobile/event/paint"
)

func Controller(r xbox.Event) error {

	e, err := r.GetEvent()
	if err != nil {
		return err
	}

	if xbox.JudgeAxis(e, xbox.CROSS_VERTICAL) {
		if e.Axes[xbox.CROSS_VERTICAL] > 0 {
			idx++
		} else {
			idx--
		}
		if idx < 0 {
			idx = 0
		} else if idx >= len(images) {
			idx = len(images) - 1
		}
		file := images[idx]
		src, err := load(file)
		if err != nil {
			return err
		}

		current = src
		win.Send(paint.Event{})
	}

	return nil
}
