// Copyright 2011 <chaishushan@gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/lazywei/go-opencv/opencv"
	//"../opencv" // can be used in forks, comment in real application
)

func init() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
}

func main() {
	filename := "matrix2.mp4"
	cap := opencv.NewFileCapture(filename)
	if cap == nil {
		panic("can not open video")
	}
	defer cap.Release()

	win := opencv.NewWindow("GoOpenCV: VideoPlayer")
	defer win.Destroy()

	fps := int(cap.GetProperty(opencv.CV_CAP_PROP_FPS))
	frames := int(cap.GetProperty(opencv.CV_CAP_PROP_FRAME_COUNT))
	stop := false

	win.SetMouseCallback(func(event, x, y, flags int) {
		if flags&opencv.CV_EVENT_LBUTTONDOWN != 0 {
			stop = !stop
			if stop {
				fmt.Printf("status: stop")
			} else {
				fmt.Printf("status: palying")
			}
		}
	})
	win.CreateTrackbar("Seek", 1, frames, func(pos int) {
		cur_pos := int(cap.GetProperty(opencv.CV_CAP_PROP_POS_FRAMES))

		if pos != cur_pos {
			cap.SetProperty(opencv.CV_CAP_PROP_POS_FRAMES, float64(pos))
			fmt.Printf("Seek to %d(%d)\n", pos, frames)
		}
	})

	for {
		if !stop {
			img := cap.QueryFrame()
			if img == nil {
				break
			}

			frame_pos := int(cap.GetProperty(opencv.CV_CAP_PROP_POS_FRAMES))
			if frame_pos >= frames {
				break
			}
			win.SetTrackbarPos("Seek", frame_pos)

			win.ShowImage(img)
			key := opencv.WaitKey(1000 / fps)
			if key == 27 {
				os.Exit(0)
			}
		} else {
			key := opencv.WaitKey(20)
			if key == 27 {
				os.Exit(0)
			}
		}
	}

	opencv.WaitKey(0)
}
