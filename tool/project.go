package tool

import (
	"fmt"
	"image"
	"image/draw"
	"os"
	"strings"

	//"github.com/ikascrew/go-opencv/opencv"
	"github.com/ikascrew/ikascrew"

	"gopkg.in/cheggaaa/pb.v1"

	"gocv.io/x/gocv"
)

type Movie struct {
	Name  string
	Type  string
	Image string
}

//const PUBLIC = ".public"
const IMAGES = "images"

const CLIENT = ".client"
const THUMB = "thumb"
const ICON = "icon"

/*
func GetPublicDir() string {
	dir := ikascrew.ProjectName()
	return dir + "/" + PUBLIC
}
*/

func GetClientDir() string {
	dir := ikascrew.ProjectName()
	return dir + "/" + CLIENT
}

func CreateProject(dir string) error {

	err := ikascrew.Load(dir)
	if err != nil {
		return fmt.Errorf("Load Project:%v", err)
	}

	//public := GetPublicDir()
	client := GetClientDir()

	//images := public + "/" + IMAGES
	thumb := client + "/" + THUMB
	icon := client + "/" + ICON

	os.RemoveAll(thumb)
	os.RemoveAll(icon)

	err = Mkdir([]string{thumb, icon})
	if err != nil {
		return fmt.Errorf("Error make directory:%v", err)
	}

	files, err := Search(dir, []string{client})
	if err != nil {
		return fmt.Errorf("Error directory search:%s", err)
	}

	movies := make([]Movie, len(files))
	bar := pb.StartNew(len(files)).Prefix("Create Thumbnail")

	//Chack app.json

	cut := 2

	for idx, f := range files {

		movie := Movie{}
		work := strings.Replace(f, dir, "", 1)

		ft := "file"

		//Size Check

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

	var ipl *gocv.Mat
	if isImage(in) {
		//ipl = opencv.LoadImage(in)
		wk := gocv.IMRead(in, gocv.IMReadColor)
		ipl = &wk
	} else {
		cap, err := gocv.VideoCaptureFile(in)
		if err != nil {
			return err
		}
		if cap == nil {
			return fmt.Errorf("New Capture Error:[%s]", in)
		}
		defer cap.Close()
		//frames := int(cap.GetProperty(opencv.CV_CAP_PROP_FRAME_COUNT))
		frames := int(cap.Get(gocv.VideoCaptureFrameCount))

		center := frames / 2
		//cap.Set(opencv.CV_CAP_PROP_POS_FRAMES, float64(center))
		cap.Set(gocv.VideoCapturePosFrames, float64(center))
		mat := gocv.NewMat()

		cap.Read(&mat)

		ipl = &mat
	}

	defer ipl.Close()

	resize := gocv.NewMat()
	defer resize.Close()

	gocv.Resize(*ipl, &resize, image.Point{}, 0.25, 0.25, gocv.InterpolationDefault)
	gocv.IMWrite(out, resize)

	return nil
}

func createThumbnail(in, out string, cut int) error {

	images := make([]gocv.Mat, cut+1)
	width := 0
	height := 0
	frame := 0

	if isImage(in) {

		ipl := gocv.IMRead(in, gocv.IMReadColor)
		defer ipl.Close()
		for idx := 0; idx < len(images); idx++ {

			images[idx] = ipl

			width += ipl.Cols()
			height = ipl.Rows()
		}

	} else {
		cap, err := gocv.VideoCaptureFile(in)
		if err != nil {
			return nil
		}
		if cap == nil {
			return fmt.Errorf("New Capture Error:[%s]", in)
		}
		defer cap.Close()

		//frames := int(cap.GetProperty(opencv.CV_CAP_PROP_FRAME_COUNT))

		frames := int(cap.Get(gocv.VideoCaptureFrameCount))
		mod := frames / cut

		for idx := 0; idx < len(images); idx++ {

			m := gocv.NewMat()
			cap.Read(&m)

			for m.Empty() {

				frame = frame - 5
				if frame < 0 {
					frame = frames
				}
				//cap.SetProperty(opencv.CV_CAP_PROP_POS_FRAMES, float64(frame))
				cap.Set(gocv.VideoCapturePosFrames, float64(frame))
				cap.Read(&m)
			}

			images[idx] = m
			if frames == 1 {
				images[1] = m
				images[2] = m
				break
			}

			width += m.Cols()
			height = m.Rows()

			frame = mod * (idx + 1)
			//cap.SetProperty(opencv.CV_CAP_PROP_POS_FRAMES, float64(frame))
			cap.Set(gocv.VideoCapturePosFrames, float64(frame))
		}
	}

	l := len(images)

	//TODO
	thumbRect := image.Rect(0, 0, 1024*l, 576)
	thumb := image.NewRGBA(thumbRect)

	left := 0
	for _, elm := range images {

		img, err := elm.ToImage()
		if err != nil {
			return err
		}

		rect := image.Rect(left, 0, elm.Cols()+left, elm.Rows())
		draw.Draw(thumb, rect, img, image.Pt(0, 0), draw.Over)

		left += elm.Cols()
	}

	mat, err := gocv.ImageToMatRGB(thumb)
	if err != nil {
		return err
	}

	resize := gocv.NewMatWithSize(height/6, left/6, gocv.MatTypeCV8UC3)
	defer resize.Close()
	gocv.Resize(mat, &resize, image.Point{}, 0.166, 0.166, gocv.InterpolationDefault)
	gocv.IMWrite(out, resize)
	return nil
}

func createPage() error {
	return nil
}
