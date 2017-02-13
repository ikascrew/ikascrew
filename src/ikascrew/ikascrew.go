package ikascrew

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
)

var project string
var videos map[string]Video

func init() {
	project = ""
	videos = make(map[string]Video)
}

func PrintVideos() {
	fmt.Printf("############################### %s\n", project)
	for k, _ := range videos {
		fmt.Println(k)
	}
	fmt.Printf("#######################################\n")
}

func GetSource(name string) (Video, error) {

	if strings.LastIndex(name, ".mp4") == -1 {
		name = name + ".mp4"
	}

	v, flg := videos[name]
	if !flg {
		return nil, fmt.Errorf("%s not Found", name)
	}

	return v, nil
}

func GetVideo(name string) (Video, error) {
	v, err := GetSource(project + "/" + name)
	if err != nil {
		v, err = GetSource(name)
	}
	return v, err
}

func SetVideo(name string, v Video) {
	videos[name] = v
}

func List() []string {

	rtn := make([]string, 0)
	for key, _ := range videos {
		d := strings.Replace(key, project+"/", "", 1)
		rtn = append(rtn, d)
	}

	sort.Strings(rtn)
	return rtn
}

func Loading(name string) error {

	if project == "" {
		project = name
		fmt.Println("Project Name:", project)
		loadPlugin()
	}

	fileInfos, err := ioutil.ReadDir(name)
	if err != nil {
		return err
	}

	for _, f := range fileInfos {
		fname := f.Name()
		if f.IsDir() {
			err = Loading(name + "/" + fname)
			if err != nil {
				return err
			}
		} else {
			idx := strings.LastIndex(fname, ".mp4")
			if idx == len(fname)-4 {
				v, err := NewFile(name + "/" + fname)
				if err != nil {
					return fmt.Errorf("Error New Video:%s", err)
				}
				videos[name+"/"+fname] = v
			}
		}
	}
	return nil
}

func loadPlugin() {
	p, _ := NewTwitter()
	videos[p.Source()] = p
	img, _ := NewImage()
	videos[img.Source()] = img
}
