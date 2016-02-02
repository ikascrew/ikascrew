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
	"github.com/lazywei/go-opencv/opencv"
)

func appMain(driver gxui.Driver) {

	filename := "matrix2.mp4"
	capt := opencv.NewFileCapture(filename)
	if capt == nil {
		panic("can not open video")
	}
	defer capt.Release()

	w := int(capt.GetProperty(opencv.CV_CAP_PROP_FRAME_WIDTH))
	h := int(capt.GetProperty(opencv.CV_CAP_PROP_FRAME_HEIGHT))

	theme := dark.CreateTheme(driver)
	imgWd := theme.CreateImage()

	window := theme.CreateWindow(w, h, "movie viewer")
	window.SetScale(1.0)
	window.AddChild(imgWd)

	rect := image.Rect(0, 0, w, h)
	rgba := image.NewRGBA(rect)

	cvImage := capt.QueryFrame()
	draw.Draw(rgba, rect, cvImage.ToImage(), image.ZP, draw.Src)

	texture := driver.CreateTexture(rgba, 1)
	imgWd.SetTexture(texture)

	window.OnClose(driver.Terminate)
}

func main() {
	gl.StartDriver(appMain)
}
