package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ikascrew/ikascrew/client"
	"github.com/ikascrew/ikascrew/server"
	"github.com/ikascrew/ikascrew/tool"
)

func main() {

	flag.Parse()

	args := flag.Args()
	l := len(args)

	if l < 1 {
		fmt.Println("Error:ikascrew [init|server|client|create(until support)]")
		os.Exit(1)
	}

	cmd := args[0]

	var project string
	if l >= 2 {
		project = args[1]
	}

	var err error

	// create mode
	// create app.json

	switch cmd {
	case "init":
		err = tool.CreateProject(project)
	case "server":
		err = server.Start(project)
	case "client":
		err = client.Start()
	default:
		err = fmt.Errorf("Error:ikascrew command[init|server|client]")
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Done!")
	os.Exit(0)
}
