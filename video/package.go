package video

import ()

type Type string

const (
	FILE  Type = "file"
	IMAGE Type = "image"
	MIC   Type = "mic"
)

func Get(p, n string) (*File, error) {
	path := p + "/" + n
	return NewFile(path)
}
