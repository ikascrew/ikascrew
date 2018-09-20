package tool

import (
	"fmt"
	"image"
	"image/jpeg"
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
			/*
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
			*/
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

	//resize := opencv.Resize(ipl, 256, 144, opencv.CV_INTER_LINEAR)
	gocv.Resize(*ipl, &resize, image.Point{}, 0.25, 0.25, gocv.InterpolationDefault)
	/*
		outFile, err := os.Create(out)
		if err != nil {
			return fmt.Errorf("Error OpenFile[%s]:%s", out, err)
		}
		defer outFile.Close()
	*/

	gocv.IMWrite(out, resize)

	/*
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
	*/
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
				//cap.SetProperty(opencv.CV_CAP_PROP_POS_FRAMES, float64(frame))
				cap.Set(gocv.VideoCapturePosFrames, float64(frame))
				cap.Read(&m)
			}

			images[idx] = m

			width += m.Cols()
			height = m.Rows()

			frame = mod * (idx + 1)
			//cap.SetProperty(opencv.CV_CAP_PROP_POS_FRAMES, float64(frame))
			cap.Set(gocv.VideoCapturePosFrames, float64(frame))
		}
	}

	//thumb := opencv.CreateImage(width, height, opencv.IPL_DEPTH_8U, 3)
	thumb := gocv.NewMatWithSize(width, height, gocv.MatTypeCV8U)
	defer thumb.Close()

	left := 0
	for _, elm := range images {

		defer elm.Close()

		//rect := gocv.NewMatWithSize(elm.Cols(), elm.Rows(), gocv.MatTypeCV8U)

		/*
			var rect gocv.Mat
			rect.Init(left, 0, elm.Cols(), elm.Rows())
		*/

		//thumb.SetROI(rect)
		//gocv.Copy(elm, thumb, nil)

		left += elm.Cols()
	}

	l := len(images)
	//thumb.ResetROI()

	resize := gocv.NewMat()
	defer resize.Close()
	//gocv.Resize(*ipl, &resize, image.Point{}, 0.25, 0.25, gocv.InterpolationDefault)
	//gocv.Resize(thumb, int(left/l/2), int(height/l/2), opencv.CV_INTER_LINEAR)
	gocv.Resize(thumb, &resize, image.Point{}, float64(left/l/2), float64(height/l/2), gocv.InterpolationDefault)

	defer resize.Close()

	img, err := resize.ToImage()
	if err != nil {
		return fmt.Errorf("Error ToImage[%s]:%s", out, err)
	}

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
