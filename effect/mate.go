package effect

import (
	"encoding/binary"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/secondarykey/go-opencv/opencv"
	"github.com/secondarykey/ikascrew"
)

type Mate struct {
	s         *Switch
	end       bool
	pressed   bool
	pressedAt time.Time
	device    *os.File
}

const (
	MATE_FILE = "/dev/input/powermate"
)

func NewMate(v ikascrew.Video, e ikascrew.Effect) (*Mate, error) {

	s, err := NewSwitch(v, e)
	if err != nil {
		return nil, err
	}

	s.number = 100

	device, err := os.OpenFile(MATE_FILE, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}

	m := Mate{
		s:       s,
		pressed: false,
		device:  device,
	}
	go m.loop()

	fmt.Printf("Mate[%p] New\n", &m)

	return &m, nil
}

func (b *Mate) Next() (*opencv.IplImage, error) {

	now, _ := b.s.now.Next()
	next, _ := b.s.video.Next()

	alpha := b.s.count / b.s.number
	opencv.AddWeighted(next, float64(alpha), now, float64(1.0-alpha), 0.0, b.s.img)

	return b.s.img, nil
}

func (b *Mate) Wait() int {
	return b.s.Wait()
}

func (b *Mate) Release() error {

	fmt.Printf("Mate[%p] Release\n", b)

	b.end = true
	return b.s.Release()
}

func (b *Mate) String() string {
	return b.s.String()
}

//github.com/awly/pmd
func (m *Mate) loop() {
	buf := make([]byte, 48)
	for {
		n, err := m.device.Read(buf)
		if err != nil {
			break
		}
		event := buf[:n]

		w := event[16:20]
		typ, _ := binary.Varint(w)

		w = event[20:24]
		val, _ := binary.Varint(w)

		m.handle(int32(typ), int32(val))
		if m.end {
			m.device.Close()
			runtime.Goexit()
		}
	}
}

const (
	// Raw event types received from the device.
	typeRot   = 1
	typePress = -1
	valLeft   = 0
	valRight  = -1
	valUp     = 0
	valDown   = -1
	// Threshold to trigger pressed rotation actions.
	// Pressed rotation actions have more impact and need to be less sensitive.
	pressedRotTicks = 2
)

func (s *Mate) handle(typ, val int32) {
	switch typ {

	case typeRot:
		switch val {
		case valRight:
			s.s.count++
			if s.pressed {
				if int(s.s.count)%pressedRotTicks == 0 {
					//trigger(evPressedRotRight)
				}
			} else {
				//trigger(evRotRight)
			}
		case valLeft:
			s.s.count--
			if s.pressed {
				if int(s.s.count)%pressedRotTicks == 0 {
					//trigger(evPressedRotLeft)
				}
			} else {
				//trigger(evRotLeft)
			}
		}
	case typePress:
		s.s.count = 0
		switch val {
		case valDown:
			s.pressed = true
			s.pressedAt = time.Now()
		case valUp:
			s.pressed = false
			if time.Since(s.pressedAt) < time.Second {
				//trigger(evClick)
			}
		}
	}
}

// Translated events from raw device events.
type event int

const (
	// Click is rapid press and release.
	evClick event = iota
	// RotRight is unpressed clockwise rotation.
	evRotRight
	// RotLeft is unpressed counter-clockwise rotation.
	evRotLeft
	// PressedRotRight is pressed clockwise rotation.
	evPressedRotRight
	// PressedRotLeft is pressed counter-clockwise rotation.
	evPressedRotLeft
)
