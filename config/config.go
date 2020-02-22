package config

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/ikascrew/ikasbox/handler"
)

type AppConfig struct {
	Width   int
	Height  int
	Default Default
}

type Default struct {
	Type string
	Name string
}

var conf *AppConfig

func init() {
	conf = nil
}

func Load(p string) (*AppConfig, error) {

	url := "http://localhost:5555/project/content/list/" + p
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	res := handler.ProjectResponse{}

	err = json.Unmarshal(byteArray, &res)
	if err != nil {
		return nil, err
	}

	def := Default{
		Type: "terminal",
		Name: "blank",
	}

	app := AppConfig{
		Width:   res.Project.Width,
		Height:  res.Project.Height,
		Default: def,
	}

	return &app, nil
}
