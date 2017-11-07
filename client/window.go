package client

import (
	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
)

type Window struct {
	window screen.Window

	list *List
	next *Next

	/*
		cursorS   int
		currentS  int
		images    []image.Image
		resourceS []string

		cursorP   int
		currentP  int
		targets   []image.Image
		resourceP []string
	*/
}

func NewWindow(t string, w, h int) (window *Window, err error) {

	window = &Window{
		window: nil,
	}

	driver.Main(func(s screen.Screen) {
		opt := &screen.NewWindowOptions{
			Title:  t,
			Width:  w,
			Height: h,
		}
		window.window, err = s.NewWindow(opt)
		if err != nil {
			return
		}
		/*
			winSize := image.Point{w, h}
			window.buffer, err = s.NewBuffer(winSize)
			if err != nil {
				return
			}
		*/

		l, err := NewList(window.window, s)
		if err != nil {
			return
		}

		n, err := NewNext(window.window, s)
		if err != nil {
			return
		}

		window.list = l
		window.next = n

	})

	return
}

func (w *Window) Release() {
	w.window.Release()
	//w.buffer.Release()
}

/*
func (w *Window) keyListener(k int) {

	fmt.Printf("[%d]\n", k)

	switch k {
	case 82:
		w.list.setCursor(-60000)
		w.Repaint()
	case 81:
		w.setCursorS(60000)
		w.Repaint()
	case 80:
		w.setCursorP(10000)
	case 79:
		w.setCursorP(-10000)
	case 40:
	}
}
*/
