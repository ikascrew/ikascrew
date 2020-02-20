package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gocv.io/x/gocv"
)

func main() {

	flag.Parse()

	args := flag.Args()

	if len(args) != 1 {
		fmt.Println("Argument FilePath")
		os.Exit(1)
	}

	err := run(args)
	if err != nil {
		fmt.Println("Argument FilePath: %v", err)
		os.Exit(1)
	}

	fmt.Println("Complete")
}

func run(args []string) error {

	path := args[0]

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return fmt.Errorf("ioutil.ReadDir(%s) :%w", path, err)
	}

	for _, file := range files {

		name := filepath.Join(path, file.Name())
		names := strings.Split(name, ".mp4")

		if len(names) != 2 {
			fmt.Println("not MP4")
			continue
		}

		newName := names[0] + "(loop).mp4"

		fmt.Printf("%s -> %s\n", name, newName)

		err := write(name, newName)
		if err != nil {
			return fmt.Errorf("write(%s) :%w", name, err)
		}
	}

	return nil
}

func write(in string, out string) error {

	w := 1280
	h := 720

	inCap, err := gocv.VideoCaptureFile(in)
	if err != nil {
		return err
	}

	frames := int(inCap.Get(gocv.VideoCaptureFrameCount))

	mixCap, err := gocv.VideoCaptureFile(in)
	if err != nil {
		return err
	}

	writer, err := gocv.VideoWriterFile(out, "X264", 29, w, h, true)
	if err != nil {
		return fmt.Errorf("error opening video writer file: %v\n", in)
	}

	inSource := gocv.NewMatWithSize(h, w, gocv.MatTypeCV8UC3)
	mixSource := gocv.NewMatWithSize(h, w, gocv.MatTypeCV8UC3)
	writeSource := gocv.NewMatWithSize(h, w, gocv.MatTypeCV8UC3)

	size := int(float64(frames) * 0.2)
	remixPos := frames - size

	inCap.Set(gocv.VideoCapturePosFrames, float64(size))
	for idx := size; idx <= frames; idx++ {

		inCap.Read(&inSource)
		if idx < remixPos {
			writer.Write(inSource)
			continue
		}

		if inSource.Empty() {
			break
		}

		mixCap.Read(&mixSource)

		alpha := float64(idx-remixPos) / float64(size)

		gocv.AddWeighted(mixSource, float64(alpha), inSource, float64(1.0-alpha), 0.0, &writeSource)

		writer.Write(writeSource)
	}
	err = writer.Close()
	if err != nil {
		return err
	}

	return nil
}
