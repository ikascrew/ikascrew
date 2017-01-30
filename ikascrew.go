package main

import (
	"fmt"
	"image/jpeg"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/secondarykey/go-opencv/opencv"
	"gopkg.in/cheggaaa/pb.v1"
)

func main() {

	dir := "setting/20161226"
	thumb := dir + "/thumb"

	err := os.MkdirAll(thumb, 0777)
	if err != nil {
		fmt.Println("Error make directory:", thumb)
		return
	}

	files, err := search(dir)
	if err != nil {
		fmt.Println("Error directory search:", err)
		return
	}

	bar := pb.StartNew(len(files)).Prefix("Create Thumbnail")
	for _, f := range files {

		work := strings.Replace(f, dir, "", 1)
		jpg := strings.Replace(work, ".mp4", ".jpg", 1)
		out := thumb + jpg

		mkIdx := strings.LastIndex(out, "/")
		tmp := string(out[:mkIdx])

		err = os.MkdirAll(tmp, 0777)

		err = createThumbnail(f, out, 5)
		if err != nil {
			fmt.Println("Error Create Thumbnail:", err)
			os.Exit(1)
		}
		bar.Increment()
	}
	bar.FinishPrint("Complate!")
	return
}

func search(d string) ([]string, error) {
	fileInfos, err := ioutil.ReadDir(d)
	if err != nil {
		return nil, fmt.Errorf("Error:Read Dir[%s]", d)
	}

	rtn := make([]string, 0)
	for _, f := range fileInfos {
		fname := f.Name()
		if f.IsDir() {
			files, err := search(d + "/" + fname)
			if err != nil {
				return nil, err
			}
			rtn = append(rtn, files...)
		} else {
			idx := strings.LastIndex(fname, ".mp4")
			if idx == len(fname)-4 {
				rtn = append(rtn, d+"/"+fname)
			}
		}
	}

	sort.Strings(rtn)
	return rtn, nil
}

func createThumbnail(in, out string, cut int) error {

	// load movie
	cap := opencv.NewFileCapture(in)
	if cap == nil {
		return fmt.Errorf("New Capture Error:[%s]", in)
	}
	defer cap.Release()

	frames := int(cap.GetProperty(opencv.CV_CAP_PROP_FRAME_COUNT))
	mod := frames / cut

	images := make([]*opencv.IplImage, cut+1)
	width := 0
	height := 0

	// get thumb
	frame := 0
	for idx := 0; idx < len(images); idx++ {

		img := cap.QueryFrame()
		for img == nil {
			frame = frame - 5
			cap.SetProperty(opencv.CV_CAP_PROP_POS_FRAMES, float64(frame))
			img = cap.QueryFrame()
		}
		images[idx] = img.Clone()

		width += img.Width()
		height = img.Height()

		frame = mod * (idx + 1)
		cap.SetProperty(opencv.CV_CAP_PROP_POS_FRAMES, float64(frame))
	}

	//create thumb

	thumb := opencv.CreateImage(width, height, opencv.IPL_DEPTH_8U, 3)
	defer thumb.Release()
	left := 0

	for _, elm := range images {

		defer elm.Release()

		var rect opencv.Rect
		rect.Init(left, 0, elm.Width(), elm.Height())

		thumb.SetROI(rect)
		opencv.Copy(elm, thumb, nil)

		left += elm.Width()
	}

	//generate thumb

	l := len(images)
	thumb.ResetROI()
	resize := opencv.Resize(thumb, int(left/l/2), int(height/l/2), opencv.CV_INTER_LINEAR)
	img := resize.ToImage()

	outFile, err := os.Create(out)
	if err != nil {
		return fmt.Errorf("Error OpenFile[%s]:%s", out, err)
	}
	defer outFile.Close()

	option := &jpeg.Options{Quality: 100}
	if err = jpeg.Encode(outFile, img, option); err != nil {
		return fmt.Errorf("Error Encode:%s", err)
	}

	return nil
}
