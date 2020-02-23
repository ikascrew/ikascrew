package client

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/golang/glog"
	"github.com/ikascrew/xbox"
)

var Y bool  //Server push
var X bool  //Server push
var A bool  //next Add
var B bool  //next Delete
var L1 bool //Server switch
var R1 bool //Server switch

func (ika *IkascrewClient) controller(e xbox.Event) error {

	if xbox.JudgeAxis(e, xbox.CROSS_VERTICAL) {
		ika.selector.list.setCursor(e.Axes[xbox.CROSS_VERTICAL] / 2)
		ika.selector.list.Push()
	}

	if xbox.JudgeAxis(e, xbox.CROSS_HORIZONTAL) {
		ika.selector.next.setCursor(e.Axes[xbox.CROSS_HORIZONTAL])
		ika.selector.next.Push()
	}

	if e.Buttons[xbox.L1] && L1 {
		ika.selector.list.maxCursor()
		/*
			L1 = false
			err := ika.callPrev()
			if err != nil {
				glog.Error("callPrev[" + err.Error() + "]")
			}
		*/
	} else if !e.Buttons[xbox.L1] {
		L1 = true
	}

	if e.Buttons[xbox.R1] && R1 {
		ika.selector.list.zeroCursor()
		/*
			R1 = false
			err := ika.callNext()
			if err != nil {
				glog.Error("callNext[" + err.Error() + "]")
			}
		*/
	} else if !e.Buttons[xbox.R1] {
		R1 = true
	}

	//Controller

	if e.Buttons[xbox.Y] && Y {
		Y = false
		res := ika.selector.next.get()
		if res != "" {

			t := "file"
			if strings.Index(res, "countdown") >= 0 {
				t = "countdown"
			}

			idx := strings.LastIndex(res, "/")
			buf := ""
			if idx != -1 {
				buf = res[idx+1:]
				buf = strings.Replace(buf, ".jpg", "", -1)
				fmt.Println(buf)
			}

			id, err := strconv.Atoi(buf)
			if err != nil {
				return fmt.Errorf("Efffect id error:[%s]", res)
			}

			err = ika.callEffect(int64(id), t)
			if err != nil {
				glog.Error("callEffect[" + err.Error() + "]")
			} else {

				//0
				setZero()

				ika.selector.next.delete()
				ika.selector.next.Push()
			}
		} else {
			glog.Error("Pusher Error: No Index")
		}
	} else if !e.Buttons[xbox.Y] {
		Y = true
	}

	if e.Buttons[xbox.X] && X {
		X = false
		res := ika.selector.list.get()
		if res != "" {

			ika.selector.player.setFile(res)
			ika.selector.player.Draw()
			ika.selector.player.Push()

		} else {
			glog.Error("Pusher Error: No Index")
		}

	} else if !e.Buttons[xbox.X] {
		X = true
	}

	if e.Buttons[xbox.B] && B {
		B = false

		res := ika.selector.list.get()
		if res != "" {
			err := ika.selector.next.add(res)
			if err != nil {
				// TODO 無視
			}
			ika.selector.next.Push()
		} else {
			glog.Error("Selector Error:" + "No Index")
		}

	} else if !e.Buttons[xbox.B] {
		B = true
	}

	if e.Buttons[xbox.A] && A {
		A = false
		err := ika.selector.next.delete()
		if err != nil {
			// TODO 無視
			glog.Error("Pusher Delete Error:", err)
		}
		ika.selector.next.Push()
	} else if !e.Buttons[xbox.A] {
		A = true
	}

	return nil
}
