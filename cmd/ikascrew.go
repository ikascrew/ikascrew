package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ikascrew/ikascrew"
	"github.com/ikascrew/ikascrew/tool"
)

func main() {

	flag.Parse()

	args := flag.Args()
	l := len(args)

	if l < 1 {
	}

	cmd := args[0]

	var project string
	var file string
	if l >= 2 {
		project = args[1]
	}

	if l >= 3 {
		file = args[2]
	}

	var err error

	// create mode
	// create app.json

	switch cmd {
	case "init":
		if project == "" {
			err = fmt.Errorf("Required:ProjectName")
		} else {
			err = tool.CreateProject(project)
		}
	case "start":
		if project == "" {
			err = fmt.Errorf("Required:ProjectName")
		} else {
			err = ikascrew.Start(project)
		}
	case "test":
		if project == "" {
			err = fmt.Errorf("Required:ProjectName")
		} else if file == "" {
			err = fmt.Errorf("Required:Filename")
		} else {
			err = ikascrew.TestMode(project, file)
		}
	default:
		err = fmt.Errorf("Error:ikascrew command[init|server|client|test]")
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Done!")
	os.Exit(0)
}
