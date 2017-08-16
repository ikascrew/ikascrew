package ikascrew

import (
	"fmt"

	"github.com/ikascrew/ikascrew/config"

	"github.com/golang/glog"
)

var project string
var Config *config.AppConfig

func init() {
	project = ""
	Config = nil
}

func ProjectName() string {
	return project
}

func Load(name string) error {

	glog.Info("Loading Project[" + name + "]")

	project = name
	conf, err := config.Load(project)
	if err != nil {
		return fmt.Errorf("Error LoadConfig[%v]", err)
	}
	Config = conf
	return nil
}
