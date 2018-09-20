package client

import (
	"github.com/ikascrew/xbox"

	"github.com/golang/glog"
)

var Y bool  //Server push
var X bool  //Server push
var A bool  //next Add
var B bool  //next Delete
var L2 bool //Server switch
var R2 bool //Server switch

// JOY_L_VERTICAL   list select
// JOY_L_HORIZONTAL next select

func (ika *IkascrewClient) controller(e xbox.Event) error {

	if xbox.JudgeAxis(e, xbox.JOY_L_VERTICAL) {
		ika.selector.list.setCursor(e.Axes[xbox.JOY_L_VERTICAL])
		ika.selector.list.Push()
	}

	if xbox.JudgeAxis(e, xbox.JOY_R_HORIZONTAL) {
		ika.selector.next.setCursor(e.Axes[xbox.JOY_R_HORIZONTAL])
		ika.selector.next.Push()
	}

	if xbox.JudgeAxis(e, xbox.L2) && L2 {
		L2 = false
		err := ika.callPrev()
		if err != nil {
			glog.Error("callPrev[" + err.Error() + "]")
		}
	} else if !xbox.JudgeAxis(e, xbox.L2) {
		L2 = true
	}

	if xbox.JudgeAxis(e, xbox.R2) && R2 {
		R2 = false
		err := ika.callNext()
		if err != nil {
			glog.Error("callNext[" + err.Error() + "]")
		}
	} else if !xbox.JudgeAxis(e, xbox.R2) {
		R2 = true
	}

	if e.Buttons[xbox.X] && X {
		X = false
		res := ika.selector.next.get()
		if res != "" {
			err := ika.callEffect(res, "file")
			if err != nil {
				glog.Error("callEffect[" + err.Error() + "]")
			} else {
				ika.selector.next.delete()
				ika.selector.next.Push()
			}
		} else {
			glog.Error("Pusher Error: No Index")
		}
	} else if !e.Buttons[xbox.X] {
		X = true
	}

	if e.Buttons[xbox.Y] && Y {
		Y = false
		res := ika.selector.list.get()
		if res != "" {

			ika.selector.player.setFile(res)
			ika.selector.player.Draw()
			ika.selector.player.Push()

		} else {
			glog.Error("Pusher Error: No Index")
		}

	} else if !e.Buttons[xbox.Y] {
		Y = true
	}

	if e.Buttons[xbox.A] && A {
		A = false

		res := ika.selector.list.get()
		if res != "" {
			err := ika.selector.next.add(res)
			if err != nil {
				// TODO 無視
				glog.Error("Pusher Add Error:", err)
			}
			ika.selector.next.Push()
		} else {
			glog.Error("Selector Error:" + "No Index")
		}

	} else if !e.Buttons[xbox.A] {
		A = true
	}

	if e.Buttons[xbox.B] && B {
		B = false
		err := ika.selector.next.delete()
		if err != nil {
			// TODO 無視
			glog.Error("Pusher Delete Error:", err)
		}
		ika.selector.next.Push()
	} else if !e.Buttons[xbox.B] {
		B = true
	}

	return nil
}
