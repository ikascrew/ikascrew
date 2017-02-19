package video

import (
	"testing"

	"github.com/secondarykey/go-opencv/opencv"
)

func TestPhrase(t *testing.T) {

	win := opencv.NewWindow("Phrase")

	tx := make([]string, 10)
	tx[0] = "表示する日本語1"
	tx[1] = "表示する日本語2"
	tx[2] = "表示する日本語3"
	tx[3] = "表示する日本語4"
	tx[4] = "表示する日本語5"
	tx[5] = "表示する日本語6"
	tx[6] = "表示する日本語7"
	tx[7] = "表示する日本語8"
	tx[8] = "I donn't like love because I love you."
	tx[9] = "@ikascrew ハッシュタグでつぶやくと画面に出てくるよ"

	//p, err := NewPhrase(tx)
	//if err != nil {
	//}

	color := opencv.NewScalar(255, 255, 204, 0)
	font := opencv.FontQt("Times", 20, color, opencv.CV_FONT_BOLD, opencv.CV_STYLE_NORMAL, 5)

	img := opencv.CreateImage(1024, 576, opencv.IPL_DEPTH_8U, 3)

	for idx := 0; idx < 10000000; idx++ {

		img.Zero()

		pos := opencv.Point{idx % 1024, 250}
		font.AddText(img, "test", pos)

		win.ShowImage(img)
		opencv.WaitKey(33)
	}

}
