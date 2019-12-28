package client

import (
	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
)

type Window struct {
	window screen.Window
	client *IkascrewClient

	list   *List
	next   *Next
	player *Player
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

		l, err := NewList(window.window, s)
		if err != nil {
			return
		}

		n, err := NewNext(window.window, s)
		if err != nil {
			return
		}

		p, err := NewPlayer(window.window, s)
		if err != nil {
			return
		}

		window.list = l
		window.next = n
		window.player = p

	})

	return
}

func (w *Window) SetClient(c *IkascrewClient) {
	w.client = c
}

func (w *Window) Release() {
	w.window.Release()
	w.list.Release()
	w.next.Release()
	w.player.Release()
}

/*
func (w *Window) keyListener(k int) {

	switch k {
	case 82:
		w.list.setCursor(-10000)
		w.list.Push()
	case 81:
		w.list.setCursor(10000)
		w.list.Push()
	case 80:
		w.next.setCursor(-20000)
		w.next.Push()
	case 79:
		w.next.setCursor(20000)
		w.next.Push()
	case 40:

		res := w.list.get()
		if res != "" {
			fmt.Println("kitayo")
			err := w.next.add(res)
			if err != nil {
				// TODO 無視
			}
			w.next.Push()

			w.player.setFile(res)
			w.player.Push()
		} else {
			fmt.Printf("res empty \n")
		}

	case 44:

		res := w.next.get()
		if res != "" {
			err := w.client.callEffect(res, "file")
			if err != nil {
			} else {
				w.next.delete()
				w.next.Push()
			}
		} else {
		}

	default:
		fmt.Printf("Not Defined[%d]\n", k)
	}
}
*/
