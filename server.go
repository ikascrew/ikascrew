package main

import (
	"fmt"
	"ikascrew/server"
	"log"
	"net/http"
	"os"
	"runtime"
)

func init() {
	fmt.Println("########################## Starting ikascrew Server")
}

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	err := server.Start("setting/20161226", "opening.mp4")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return
}
