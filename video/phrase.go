package video

import (
	"math/rand"
	"time"

	"github.com/secondarykey/go-opencv/opencv"
)

var fonts []*opencv.Font

func init() {

	color1 := opencv.NewScalar(255, 255, 204, 0)
	color2 := opencv.NewScalar(0, 255, 51, 0)
	color3 := opencv.NewScalar(0, 102, 255, 0)
	color4 := opencv.NewScalar(204, 255, 102, 0)
	color5 := opencv.NewScalar(201, 51, 0, 0)

	fonts = make([]*opencv.Font, 5)

	fonts[0] = opencv.FontQt("Times", 60, color1, opencv.CV_FONT_BOLD, opencv.CV_STYLE_NORMAL, 5)
	fonts[1] = opencv.FontQt("Times", 60, color2, opencv.CV_FONT_BOLD, opencv.CV_STYLE_NORMAL, 5)
	fonts[2] = opencv.FontQt("Times", 60, color3, opencv.CV_FONT_BOLD, opencv.CV_STYLE_NORMAL, 5)
	fonts[3] = opencv.FontQt("Times", 60, color4, opencv.CV_FONT_BOLD, opencv.CV_STYLE_NORMAL, 5)
	fonts[4] = opencv.FontQt("Times", 60, color5, opencv.CV_FONT_BOLD, opencv.CV_STYLE_NORMAL, 5)

	rand.Seed(time.Now().UnixNano())
}

type Phrase struct {
	text []string
	num  int
	now  []line
	dst  *opencv.IplImage
}

type line struct {
	text string
	x    int
	y    int
	rate int
	font *opencv.Font
}

var bg *opencv.IplImage

func NewPhrase(texts []string) (*Phrase, error) {

	p := Phrase{
		num: 5,
		now: make([]line, 5),
	}

	p.text = make([]string, len(texts))
	copy(p.text, texts)

	//p.dst = opencv.CreateImage(1024, 576, opencv.IPL_DEPTH_8U, 3)
	bg = opencv.LoadImage("projects/20170213/utopia.jpg")

	return &p, nil
}

func (p *Phrase) initialize() {

	var wk *opencv.IplImage
	wk = p.dst
	defer wk.Release()
	p.dst = bg.Clone()

	for idx, _ := range p.now {

		p.now[idx].x -= p.now[idx].rate
		minx := len(p.now[idx].text) * -20

		if p.now[idx].x < minx || p.now[idx].text == "" {
			rt := rand.Intn(len(p.text))
			p.now[idx].text = p.text[rt]

			p.now[idx].x = 1024

			ry := rand.Intn(500)
			p.now[idx].y = ry + 40
			rr := rand.Intn(7)
			p.now[idx].rate = rr + 3
			rf := rand.Intn(5)
			p.now[idx].font = fonts[rf]
		}
	}
	return
}

func (p *Phrase) Next() *opencv.IplImage {

	p.initialize()
	for _, elm := range p.now {
		pos := opencv.Point{elm.x, elm.y}
		elm.font.AddText(p.dst, elm.text, pos)
	}
	return p.dst
}

func (v *Phrase) Wait() int {
	return 33
}

func (v *Phrase) Size() int {
	return 100
}

func (v *Phrase) Current() int {
	return 30
}

func (v *Phrase) Set(f int) {
}

func (v *Phrase) Reload() {
	v.Set(0)
}

func (v *Phrase) Release() {
	v.dst.Release()
}

func (v *Phrase) Source() string {
	return "ikascrew_Phrase"
}
