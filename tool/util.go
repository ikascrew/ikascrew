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

	sort.Slice(rtn, func(i, j int) bool {
		pathA := rtn[i]
		pathB := rtn[j]

		slashA := len(strings.Split(pathA, "/"))
		slashB := len(strings.Split(pathB, "/"))
		if slashA != slashB {
			return slashA < slashB
		}
		return pathA < pathB
	})
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

func Mkdir(dirs ...string) error {
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
