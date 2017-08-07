package tool

import (
	"fmt"
	"image/jpeg"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/ikascrew/go-opencv/opencv"
	"gopkg.in/cheggaaa/pb.v1"
)

type Movie struct {
	Name  string
	Type  string
	Image string
}

const TMP = ".tmp"
const THUMB = "thumb"
const ICON = "icon"

func CreateProject(dir string) error {

	tmp := dir + "/" + TMP
	thumb := tmp + "/" + THUMB
	icon := tmp + "/" + ICON

	err := os.MkdirAll(thumb, 0777)
	if err != nil {
		return fmt.Errorf("Error make directory:%s", thumb)
	}

	err = os.MkdirAll(icon, 0777)
	if err != nil {
		return fmt.Errorf("Error make directory:%s", icon)
	}

	files, err := search(dir)
	if err != nil {
		return fmt.Errorf("Error directory search:%s", err)
	}

	movies := make([]Movie, len(files))
	bar := pb.StartNew(len(files)).Prefix("Create Thumbnail")

	cut := 2

	for idx, f := range files {

		movie := Movie{}
		work := strings.Replace(f, dir, "", 1)

		ft := "file"

		if isImage(work) {

			ft = "image"
			movie.Image = THUMB + work

			out := thumb + work
			mkIdx := strings.LastIndex(out, "/")
			d := string(out[:mkIdx])
			err = os.MkdirAll(d, 0777)

			err = createThumbnail(f, out, cut)
			if err != nil {
				return fmt.Errorf("Error Create Thumbnail:%s", err)
			}

			out = icon + work
			mkIdx = strings.LastIndex(out, "/")
			tmp = string(out[:mkIdx])
			err = os.MkdirAll(tmp, 0777)

			err = createIcon(f, out)
			if err != nil {
				return fmt.Errorf("Error Create Icon:%s", err)
			}

		} else {

			jpg := strings.Replace(work, ".mp4", ".jpg", 1)
			movie.Image = THUMB + jpg
			out := thumb + jpg

			mkIdx := strings.LastIndex(out, "/")
			tmp := string(out[:mkIdx])

			err = os.MkdirAll(tmp, 0777)
			err = createThumbnail(f, out, cut)
			if err != nil {
				return fmt.Errorf("Error Create Thumbnail:%s", err)
			}

			out = icon + jpg
			mkIdx = strings.LastIndex(out, "/")
			tmp = string(out[:mkIdx])
			err = os.MkdirAll(tmp, 0777)

			err = createIcon(f, out)
			if err != nil {
				return fmt.Errorf("Error Create Icon:%s", err)
			}

		}

		movie.Name = string(work[1:])
		movie.Type = string(ft)

		movies[idx] = movie
		bar.Increment()
	}
	bar.FinishPrint("Thumbnail Completion")

	tw := Movie{
		Name: "ikascrew_microphone",
		Type: "mic",
	}
	movies = append(movies, tw)
	bar.FinishPrint("Controller Completion")

	return nil
}

func isImage(f string) bool {
	if strings.LastIndex(f, ".jpg") == len(f)-4 ||
		strings.LastIndex(f, ".png") == len(f)-4 {
		return true
	}
	return false
}

func search(d string) ([]string, error) {

	fileInfos, err := ioutil.ReadDir(d)
	if err != nil {
		return nil, fmt.Errorf("Error:Read Dir[%s]", d)
	}

	rtn := make([]string, 0)
	if strings.Index(d, TMP) != -1 {
		return rtn, nil
	}

	for _, f := range fileInfos {
		fname := f.Name()
		if f.IsDir() {
			files, err := search(d + "/" + fname)
			if err != nil {
				return nil, err
			}
			rtn = append(rtn, files...)
		} else {
			midx := strings.LastIndex(fname, ".mp4")
			jidx := strings.LastIndex(fname, ".jpg")
			pidx := strings.LastIndex(fname, ".png")
			if midx == len(fname)-4 ||
				jidx == len(fname)-4 ||
				pidx == len(fname)-4 {
				rtn = append(rtn, d+"/"+fname)
			}
		}
	}

	sort.Strings(rtn)
	return rtn, nil
}

func createIcon(in, out string) error {

	var ipl *opencv.IplImage
	if isImage(in) {
		ipl = opencv.LoadImage(in)
	} else {
		cap := opencv.NewFileCapture(in)
		if cap == nil {
			return fmt.Errorf("New Capture Error:[%s]", in)
		}
		defer cap.Release()
		frames := int(cap.GetProperty(opencv.CV_CAP_PROP_FRAME_COUNT))

		center := frames / 2
		cap.SetProperty(opencv.CV_CAP_PROP_POS_FRAMES, float64(center))
		img := cap.QueryFrame()
		ipl = img.Clone()
	}

	defer ipl.Release()
	resize := opencv.Resize(ipl, 512, 288, opencv.CV_INTER_LINEAR)
	defer resize.Release()

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

func createThumbnail(in, out string, cut int) error {

	images := make([]*opencv.IplImage, cut+1)
	width := 0
	height := 0
	frame := 0

	if isImage(in) {

		ipl := opencv.LoadImage(in)
		defer ipl.Release()
		for idx := 0; idx < len(images); idx++ {
			images[idx] = ipl.Clone()
			width += ipl.Width()
			height = ipl.Height()
		}

	} else {
		cap := opencv.NewFileCapture(in)
		if cap == nil {
			return fmt.Errorf("New Capture Error:[%s]", in)
		}
		defer cap.Release()

		frames := int(cap.GetProperty(opencv.CV_CAP_PROP_FRAME_COUNT))
		mod := frames / cut

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
	}

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

	l := len(images)
	thumb.ResetROI()
	resize := opencv.Resize(thumb, int(left/l/2), int(height/l/2), opencv.CV_INTER_LINEAR)
	defer resize.Release()

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

func copyFile(src, dst string) error {
	// read the whole file at once
	b, err := ioutil.ReadFile(src)
	if err != nil {
		return fmt.Errorf("Error:Read File %s", src)
	}

	err = ioutil.WriteFile(dst, b, 0644)
	if err != nil {
		return fmt.Errorf("Error:Write File %s", dst)
	}
	return nil
}
