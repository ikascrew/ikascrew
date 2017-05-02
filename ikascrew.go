package ikascrew

import (
	"fmt"
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
	return nil
}
