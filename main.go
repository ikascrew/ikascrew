// Copyright 2011 <chaishushan@gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"

	"github.com/secondarykey/go-opencv/opencv"
)

func main() {

	filename := "/home/secondarykey/Videos/matrix2.mp4"
	capt := opencv.NewFileCapture(filename)
	if capt == nil {
		panic("can not open video")
	}
	defer capt.Release()

	fps := int(capt.GetProperty(opencv.CV_CAP_PROP_FPS))
	frames := int(capt.GetProperty(opencv.CV_CAP_PROP_FRAME_COUNT))
	width := capt.GetProperty(opencv.CV_CAP_PROP_FRAME_WIDTH)
	height := capt.GetProperty(opencv.CV_CAP_PROP_FRAME_HEIGHT)

	stop := false

	fourcc := opencv.FOURCC(int8('M'), int8('J'), int8('P'), int8('G'))
	//fourcc := opencv.FOURCC(int8('P'), int8('I'), int8('M'), int8('1'))
	//fourcc := opencv.FOURCC(int8('X'), int8('V'), int8('I'), int8('D'))

	writer := opencv.NewVideoWriter("matrix2.avi", int(fourcc), float32(fps), int(width), int(height), 1)
	if writer == nil {
		panic("New VideoWriter Error")
	}
	defer writer.Release()

	imgAfter := capt.QueryFrame()
	if imgAfter == nil {
		panic("error QueryFrame")
	}
	//dst := opencv.CreateImage(int(width), int(height), opencv.IPL_DEPTH_8U, 3)

	dst := imgAfter.Clone()

	for {
		if !stop {

			img := capt.QueryFrame()
			frame_pos := int(capt.GetProperty(opencv.CV_CAP_PROP_POS_FRAMES))
			if frame_pos >= frames {
				break
			}
			imgAfter = img.Clone()

			opencv.AddWeighted(imgAfter, float64(0.3), img, float64(0.7), 0.0, dst)

			writer.WriteFrame(dst)

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
	//opencv.WaitKey(0)
}
