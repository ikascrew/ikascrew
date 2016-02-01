// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"

	"github.com/google/gxui"
	"github.com/google/gxui/drivers/gl"
	"github.com/google/gxui/themes/dark"
	"github.com/secondarykey/go-opencv/opencv"
)

func appMain(driver gxui.Driver) {

	filename := "/home/secondarykey/Videos/matrix2.mp4"
	capt := opencv.NewFileCapture(filename)
	if capt == nil {
		panic("can not open video")
	}
	defer capt.Release()

	//
	//fps := int(capt.GetProperty(opencv.CV_CAP_PROP_FPS))
	//frames := int(capt.GetProperty(opencv.CV_CAP_PROP_FRAME_COUNT))
	w := int(capt.GetProperty(opencv.CV_CAP_PROP_FRAME_WIDTH))
	h := int(capt.GetProperty(opencv.CV_CAP_PROP_FRAME_HEIGHT))

	theme := dark.CreateTheme(driver)
	imgWd := theme.CreateImage()

	window := theme.CreateWindow(w, h, "movie viewer")
	window.SetScale(1.0)
	window.AddChild(imgWd)

	rect := image.Rect(0, 0, w, h)
	rgba := image.NewRGBA(rect)

	window.OnClose(driver.Terminate)

	// loop  ReDraw

	cvImage := capt.QueryFrame()
	frame_pos := int(capt.GetProperty(opencv.CV_CAP_PROP_POS_FRAMES))
	if frame_pos >= frames {
		break
	}

	/*
		type Image interface {
		    image.Image
		    Set(x, y int, c color.Color)
		}
	*/

	/*
	   func Draw(dst Image, r image.Rectangle, src image.Image, sp image.Point, op Op)
	   -> draw.Draw(rgba, source.Bounds(), source, image.ZP, draw.Src)
	*/

	draw.Draw(rgba, rect, cvImage.ToImage(), image.ZP, draw.Src)

	texture := driver.CreateTexture(rgba, 1)
	imgWd.SetTexture(texture)

}

func main() {
	gl.StartDriver(appMain)
}
