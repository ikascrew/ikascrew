package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"strings"

	"ikascrew"
)

var sc = bufio.NewScanner(os.Stdin)

func init() {
	fmt.Println("########################## Starting ikascrew system")
}

func main() {

	//max
	runtime.GOMAXPROCS(runtime.NumCPU())

	//debug
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	//arg directory

	err := ikascrew.Run("setting")
	if err != nil {
		os.Exit(-1)
	}

	setMovieWindow("0", "1024x576.mp4")
	//setMovieWindow("1", "default.mp4")

	fmt.Println("###################################################")
	source := "0"

	for {

		fmt.Print(source + " > ")
		sc.Scan()
		cmd := sc.Text()
		fmt.Println(cmd)

		cmds := strings.Split(cmd, " ")

		switch cmds[0] {
		case "set":
			if source != "" {
				m := cmds[1]
				setMovieWindow(source, m)
			}
		case "source":
			source = cmds[1]
		case "q":
			ikascrew.Release()
			break
		default:
		}

	}
}

func setMovieWindow(w, m string) {

	win := ikascrew.GetWindow(w)
	if win == nil {
		fmt.Println("Error GetWindow:" + w)
		return
	}

	video := ikascrew.NewBlendVideo("snow.mp4", "sumi.mp4")
	go win.Play(video)

}
