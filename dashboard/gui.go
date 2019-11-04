package dashboard

import (
	"github.com/Alethio/memento/core"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("module", "gui")

type Config struct {
	Port          string
	ConfigEnabled bool
}

type Dashboard struct {
	config Config
	engine *gin.Engine

	core *core.Core
}

func New(core *core.Core, config Config) *Dashboard {
	return &Dashboard{
		config: config,
		core:   core,
	}
}

func (d *Dashboard) Run() {
	d.engine = gin.Default()
	d.setRoutes()

	err := d.engine.Run(":" + d.config.Port)
	if err != nil {
		log.Fatal(err)
	}
}

func (d *Dashboard) Close() {
}
