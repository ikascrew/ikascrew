package ikascrew

import (
	//"github.com/ikascrew/xbox"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"golang.org/x/exp/shiny/screen"

	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
)

var current image.Image
var idx int
var win screen.Window
var images []string
var resources []string

var X bool
var A bool

func init() {
}

func exchange(dir string) {

	work := dir + "/.public/thumb"

	paths, err := search(work)
	if err != nil {
		panic(err)
	}
	images = paths

	resources = make([]string, len(paths))
	for idx, path := range images {

		jpg := strings.Replace(path, work, "", -1)
		mpg := strings.Replace(jpg, ".jpg", ".mp4", -1)
		path := mpg
		resources[idx] = path
	}

}

func display(dir string) error {

	return nil
}

func load(filename string) (image.Image, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	m, _, err := image.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("could not decode %s: %v", filename, err)
	}
	return m, nil
}

func draw(m *image.RGBA, img image.Image) {

	if img == nil {
		return
	}

	b := m.Bounds()
	lox := b.Min.X
	loy := b.Min.Y
	hix := b.Max.X
	hiy := b.Max.Y

	for y := loy; y < hiy; y++ {
		for x := lox; x < hix; x++ {
			m.Set(x, y, img.At(x, y))
		}
	}
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
