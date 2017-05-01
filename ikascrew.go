package ikascrew

import (
	"fmt"
	"io/ioutil"
	"strings"
)

var project string

func init() {
	project = ""
}

func ProjectName() string {
	return project
}

func Loading(name string) error {

	if project == "" {
		project = name
		fmt.Println("Project Name:", project)
	}

	fileInfos, err := ioutil.ReadDir(name)
	if err != nil {
		return err
	}

	//

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
			}
		}
	}
	return nil
}
