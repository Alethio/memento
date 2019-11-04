package gui

import (
	"github.com/Alethio/memento/core"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("module", "gui")

type Config struct {
	Port string
}

type GUI struct {
	config Config
	engine *gin.Engine

	core *core.Core
}

func New(core *core.Core, config Config) *GUI {
	return &GUI{
		config: config,
		core:   core,
	}
}

func (g *GUI) Run() {
	g.engine = gin.Default()
	g.setRoutes()

	err := g.engine.Run(":" + g.config.Port)
	if err != nil {
		log.Fatal(err)
	}
}

func (g *GUI) Close() {
}
