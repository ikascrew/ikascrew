package video

import (
	"fmt"
	"image"
	"image/color"

	"github.com/ikascrew/ikascrew"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"gocv.io/x/gocv"
)

func init() {
}

//var Target = time.Date(2019, time.December, 31, 9, 2, 0, 0, loc)
type Terminal struct {
	frames int
	name   string
	source *gocv.Mat

	lines []string

	now int
	max int
}

func NewTerminal(file string) (*Terminal, error) {

	f := Terminal{
		name: file,
	}
	var err error

	source := gocv.NewMatWithSize(ikascrew.Config.Height, ikascrew.Config.Width, gocv.MatTypeCV8UC3)
	f.source = &source

	f.frames = 200

	cs, err := cpu.Info()
	cpuLine := make([]string, 0)
	if err == nil {
		c := cs[0]
		cpuLine = append(cpuLine, fmt.Sprintf("    CPU -> %s x %d x %d", c.ModelName, c.Cores, len(cs)))
	} else {
		cpuLine = append(cpuLine, fmt.Sprintf("    CPU Error :%s ", err.Error()))

	}

	memLine := make([]string, 0)
	m, err := mem.VirtualMemory()
	if err == nil {
		// structが返ってきます。
		memLine = append(memLine, fmt.Sprintf("    Mem:Total: %v, Free:%v", m.Total, m.Free))
	} else {
		memLine = append(memLine, fmt.Sprintf("    Mem Error :%s ", err.Error()))
	}

	dispLine := make([]string, 0)
	dispLine = append(dispLine, fmt.Sprintf("    DISPLAY:%d x %d", ikascrew.Config.Width, ikascrew.Config.Height))

	f.lines = make([]string, 8+len(cpuLine)+len(memLine)+len(dispLine))
	//CPU
	//MEM
	f.lines[0] = "I am ikascrew."
	f.lines[1] = "I am a program born to transform \"VJ System\"."
	f.lines[2] = ""
	f.lines[3] = "Today's system:"

	idx := 4
	for _, line := range cpuLine {
		f.lines[idx] = line
		idx++
	}

	for _, line := range memLine {
		f.lines[idx] = line
		idx++
	}

	for _, line := range dispLine {
		f.lines[idx] = line
		idx++
	}

	f.lines[idx] = ""

	f.lines[idx+1] = "I am a ready."
	f.lines[idx+2] = "When you're ready?"
	f.lines[idx+3] = "Let's get started!"

	f.now = 0
	f.max = 0

	for _, line := range f.lines {
		f.max += len(line)
	}

	return &f, nil
}

func (v *Terminal) Next() (*gocv.Mat, error) {

	left := 20
	height := 30
	fps := 4

	//終了文字数
	n := v.now / fps

	v.source.Close()

	newV := gocv.NewMatWithSize(ikascrew.Config.Height, ikascrew.Config.Width, gocv.MatTypeCV8UC3)
	v.source = &newV

	for idx, line := range v.lines {

		buf := line

		charnum := len(line)

		if n < charnum {
			buf = line[0:n] + "|"
		}

		n -= len(line)

		gocv.PutText(v.source, buf, image.Pt(left, (idx+1)*height),
			gocv.FontHersheyComplexSmall, 1.0, color.RGBA{0, 255, 0, 0}, 2)

		//calet
		if n <= 0 {
			break
		}
	}

	v.now++
	return v.source, nil
}

func (v *Terminal) Set(f int) {
}

func (v *Terminal) Current() int {
	return 1
}

func (v *Terminal) Size() int {
	return v.frames
}

func (v *Terminal) Source() string {
	return v.name
}

func (v *Terminal) Release() error {
	v.source.Close()
	return nil
}
