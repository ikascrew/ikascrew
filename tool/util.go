package tool

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"image"
	_ "image/jpeg"
	_ "image/png"

	"github.com/ikascrew/go-opencv/opencv"
)

func Search(d string, ignore []string) ([]string, error) {

	fileInfos, err := ioutil.ReadDir(d)
	if err != nil {
		return nil, fmt.Errorf("Error:Read Dir[%s]", d)
	}

	rtn := make([]string, 0)

	if ignore != nil {
		for _, ig := range ignore {
			if strings.Index(d, ig) != -1 {
				return rtn, nil
			}
		}
	}

	for _, f := range fileInfos {
		fname := f.Name()
		if f.IsDir() {
			files, err := Search(d+"/"+fname, ignore)
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

func CopyFile(src, dst string) error {
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

func Mkdir(dirs []string) error {
	for _, dir := range dirs {
		err := os.MkdirAll(dir, 0777)
		if err != nil {
			return fmt.Errorf("Make directory Error:[%s]", dir)
		}
	}
	return nil
}

func LoadImage(f string) (image.Image, error) {
	d, err := os.Open(f)
	if err != nil {
		return nil, err
	}
	defer d.Close()
	m, _, err := image.Decode(d)
	if err != nil {
		return nil, fmt.Errorf("could not decode %s: %v", f, err)
	}
	return m, nil
}

func CreateMovie(f string) error {

	// exist file
	_, err := os.Stat(f)
	if err != nil {
		return err
	}

	//Change name
	idx := strings.LastIndex(f, ".")
	if idx == -1 {
		return fmt.Errorf("Error Unknown ext[%s]", f)
	}

	mp4 := string(f[:idx]) + ".mp4"
	_, err = os.Stat(mp4)
	if err == nil {
		return fmt.Errorf("Error [%s] is exist", mp4)
	}

	ipl := opencv.LoadImage(f)
	if ipl == nil {
		return fmt.Errorf("Error LoadImage[%s]", f)
	}
	defer ipl.Release()

	//mp4
	//fourcc := opencv.FOURCC('M', 'P', '4', '3')
	//fourcc := opencv.FOURCC('M', 'P', '4', '2')
	fourcc := opencv.FOURCC('X', '2', '6', '4')
	fps := 30
	w := opencv.NewVideoWriter(mp4, int(fourcc), float32(fps), ipl.Width(), ipl.Height(), 1)
	for i := 0; i < fps; i++ {
		w.WriteFrame(ipl.Clone())
	}
	w.Release()

	fmt.Printf("Output movie[%s]\n", mp4)

	return nil
}
