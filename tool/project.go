package tool

import (
	"fmt"
	"image/jpeg"
	"os"
	"strings"

	"github.com/ikascrew/go-opencv/opencv"
	"github.com/ikascrew/ikascrew"

	"gopkg.in/cheggaaa/pb.v1"
)

type Movie struct {
	Name  string
	Type  string
	Image string
}

const PUBLIC = ".public"
const IMAGES = "images"

const CLIENT = ".client"
const THUMB = "thumb"
const ICON = "icon"

func GetPublicDir() string {
	dir := ikascrew.ProjectName()
	return dir + "/" + PUBLIC
}

func GetClientDir() string {
	dir := ikascrew.ProjectName()
	return dir + "/" + CLIENT
}

func CreateProject(dir string) error {

	err := ikascrew.Load(dir)
	if err != nil {
		return fmt.Errorf("Load Project:%v", err)
	}

	public := GetPublicDir()
	client := GetClientDir()

	images := public + "/" + IMAGES
	thumb := client + "/" + THUMB
	icon := client + "/" + ICON

	os.RemoveAll(images)
	os.RemoveAll(thumb)
	os.RemoveAll(icon)

	err = Mkdir([]string{thumb, icon, images})
	if err != nil {
		return fmt.Errorf("Error make directory:%v", err)
	}

	files, err := Search(dir, []string{public, client})
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
			tmp := string(out[:mkIdx])
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
	resize := opencv.Resize(ipl, 256, 144, opencv.CV_INTER_LINEAR)
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

func createPage() error {
	return nil
}
