package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ikascrew/ikasbox/handler"
)

type AppConfig struct {
	Width   int
	Height  int
	Default Default

	Contents map[int]*Content
}

type Content struct {
	ContentID int
	Name      string
	Path      string
}

type Default struct {
	Type string
	Name string
}

var conf *AppConfig

func init() {
	conf = nil
}

func Load(p int) (*AppConfig, error) {

	url := fmt.Sprintf("http://localhost:5555/project/content/list/%d", p)
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

	app.Contents = make(map[int]*Content)

	for _, elm := range res.Contents {
		con := Content{}
		con.Name = elm.Name
		con.Path = elm.Path
		con.ContentID = elm.ContentID

		app.Contents[elm.ID] = &con
	}

	return &app, nil
}
