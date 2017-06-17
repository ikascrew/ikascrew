package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/secondarykey/ikascrew"
)

type AppConfig struct {
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

func Get() (*AppConfig, error) {

	if conf == nil {
		wk, err := load()
		if err != nil {
			return nil, err
		}
		conf = wk
	}
	return conf, nil
}

func load() (*AppConfig, error) {

	raw, err := ioutil.ReadFile(ikascrew.ProjectName() + "/app.json")
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
