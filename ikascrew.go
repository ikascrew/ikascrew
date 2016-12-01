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

	runtime.GOMAXPROCS(runtime.NumCPU())
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	err := ikascrew.Run("setting")
	if err != nil {
		os.Exit(-1)
	}

	fmt.Println("###################################################")
	source := "ikascrew"

	for {

		fmt.Print(source + " > ")
		sc.Scan()
		cmd := sc.Text()

		cmds := strings.Split(cmd, " ")

		switch cmds[0] {
		case "create":
			if len(cmds) == 3 {
				w := cmds[1]
				v := cmds[2]
				err := ikascrew.PlaySubWindow(w, v)
				if err != nil {
					fmt.Printf("Error:create command:%s\n", err)
				} else {
					source = w
				}
			} else {
				fmt.Println("Error:create command arg 2[create {window} {movie}]")
			}

		case "q":
			fmt.Println("\nByeBye?[Y/n]")
			sc.Scan()
			yN := sc.Text()
			if yN == "Y" {
				fmt.Println("Bye!")
				ikascrew.Release()
				os.Exit(0)
			}
		default:
		}

	}
}
