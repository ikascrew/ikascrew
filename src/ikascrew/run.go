package ikascrew

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"
	//"github.com/secondarykey/go-opencv/opencv"
)

var videos map[string]Queue
var windows map[string]Window

func init() {
	videos = make(map[string]Queue)
	windows = make(map[string]Window)
}

var project string

func Run(dir string) error {

	project = dir + "/data"
	err := loading(project)
	if err != nil {
		return fmt.Errorf("Error Loading directory:%s", err)
	}

	for k, _ := range videos {
		fmt.Println(k)
	}

	main := "ikascrew"
	win := NewMainWindow(main)
	windows[main] = win

	fmt.Println("OpenCL Init wait...")
	time.Sleep(5000 * time.Millisecond)

	video := NewMainVideo("snow.mp4")

	go func() {
		win.Play(video)
	}()

	return nil
}

func Release() {

	for _, v := range videos {
		v.Release()
	}

	for _, win := range windows {
		win.Destroy()
	}
}

func GetVideo(name string) (Queue, error) {

	v, flg := videos[project+"/"+name]
	if !flg {
		return nil, fmt.Errorf("%s not Found", name)
	}

	return v, nil
}

func GetWindow(name string) Window {
	return windows[name]
}

func loading(name string) error {

	fileInfos, err := ioutil.ReadDir(name)
	if err != nil {
		return err
	}

	for _, f := range fileInfos {

		fname := f.Name()
		if f.IsDir() {
			return loading(name + "/" + fname)
		} else {
			idx := strings.LastIndex(fname, ".mp4")
			if idx == len(fname)-4 {
				v, err := NewVideo(name + "/" + fname)
				if err != nil {
					return fmt.Errorf("Error ViewVide:%s", err)
				}
				videos[name+"/"+fname] = v
			}
		}
	}
	return nil
}
