package config

import (
	"encoding/json"
	"io/ioutil"
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

	raw, err := ioutil.ReadFile(p + "/app.json")
	if err != nil {
		return nil, err
	}
	app := &AppConfig{}

	err = json.Unmarshal(raw, app)
	if err != nil {
		return nil, err
	}
	return app, nil
}
