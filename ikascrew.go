package ikascrew

import (
	"fmt"

	"github.com/ikascrew/ikascrew/config"
)

var project int
var Config *config.AppConfig

func init() {
	project = 0
	Config = nil
}

func ProjectID() int {
	return project
}

func Load(id int64) error {

	var err error
	project = int(id)
	if err != nil {
		return fmt.Errorf("project id miss: %s[%v]", id, err)
	}

	conf, err := config.Load(project)
	if err != nil {
		return fmt.Errorf("Error LoadConfig[%v]", err)
	}
	Config = conf
	return nil
}
