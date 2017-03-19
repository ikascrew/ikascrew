package tool

import (
	"fmt"
	"html/template"
	"image/jpeg"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/secondarykey/go-opencv/opencv"
	"gopkg.in/cheggaaa/pb.v1"
)

type Movie struct {
	Name  string
	Image string
}

func CreateProject(dir string) error {

	public := dir + "/.public"
	thumb := public + "/thumb"
	images := public + "/images"

	err := os.MkdirAll(thumb, 0777)
	if err != nil {
		return fmt.Errorf("Error make directory:%s", thumb)
	}
	err = os.MkdirAll(images, 0777)
	if err != nil {
		return fmt.Errorf("Error make directory:%s", images)
	}

	files, err := search(dir)
	if err != nil {
		return fmt.Errorf("Error directory search:%s", err)
	}

	movies := make([]Movie, len(files))

	bar := pb.StartNew(len(files)).Prefix("Create Thumbnail")
	for idx, f := range files {

		movie := Movie{}

		work := strings.Replace(f, dir, "", 1)

		jpg := strings.Replace(work, ".mp4", ".jpg", 1)

		movie.Name = string(work[1:])
		movie.Image = "thumb" + jpg

		out := thumb + jpg

		mkIdx := strings.LastIndex(out, "/")
		tmp := string(out[:mkIdx])

		err = os.MkdirAll(tmp, 0777)

		err = createThumbnail(f, out, 5)
		if err != nil {
			return fmt.Errorf("Error Create Thumbnail:%s", err)
		}

		movies[idx] = movie
		bar.Increment()
	}
	bar.FinishPrint("Thumbnail Completion")

	tw := Movie{
		Name: "_ikascrew_Twitter.mp4",
	}
	movies = append(movies, tw)

	bar = pb.StartNew(4).Prefix("Generate Controll page")
	tmpl, err := template.ParseFiles("templates/index.tmpl")
	if err != nil {
		return fmt.Errorf("Error Create Template:%s", err)
	}

	index, err := os.Create(public + "/index.html")
	if err != nil {
		return fmt.Errorf("Error Create index:%s", err)
	}

	err = tmpl.Execute(index, movies)
	if err != nil {
		return fmt.Errorf("Error Create Index:%s", err)
	}
	bar.Increment()

	err = copyFile("templates/styles.css", public+"/styles.css")
	if err != nil {
		return fmt.Errorf("Error Copy:%s", err)
	}
	bar.Increment()

	err = copyFile("templates/images/logo.png", images+"/logo.png")
	if err != nil {
		return fmt.Errorf("Error Copy:%s", err)
	}
	bar.Increment()

	err = copyFile("templates/jquery-3.1.1.min.js", public+"/jquery-3.1.1.min.js")
	if err != nil {
		return fmt.Errorf("Error Copy:%s", err)
	}
	bar.Increment()

	bar.FinishPrint("Controller Completion")

	return nil
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
