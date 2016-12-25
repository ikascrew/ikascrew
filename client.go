package main

import (
	"fmt"
	"ikascrew/client"
	"os"
)

func main() {
	err := client.Start("setting/20161226")
	if err != nil {
		fmt.Println("Error Client Start:", err)
		os.Exit(1)
	}
	os.Exit(0)
}
