package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ikascrew/ikascrew/client"
	"github.com/ikascrew/ikascrew/server"
	"github.com/ikascrew/ikascrew/tool"
)

const VERSION = "0.0.0"

func main() {

	flag.Parse()
	args := flag.Args()

	l := len(args)
	fmt.Println(args)
	if l < 1 {
		os.Exit(1)
	}

	cmd := args[0]
	var project string
	if l >= 2 {
		project = args[1]
	}

	var err error
	switch cmd {
	case "init":
		if project == "" {
			err = fmt.Errorf("Required:ProjectName")
		} else {
			err = tool.CreateProject(project)
		}
	case "server":
		if project == "" {
			err = fmt.Errorf("Required:ProjectName")
		} else {
			err = server.Start(project)
		}
	case "client":
		err = client.Start()
	case "tool":
		if project == "" {
			err = fmt.Errorf("Required:FileName")
		} else {
			//err = tool.CreateMovie(project)
		}

	case "version":
		fmt.Println("ikascrew " + VERSION)
	default:
		err = fmt.Errorf("Error:ikascrew command[server|client]")
	}

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	os.Exit(0)
}
