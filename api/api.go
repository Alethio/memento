package api

import (
	"database/sql"

	"github.com/Alethio/memento/metrics"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("module", "api")

type Config struct {
	Port           string
	DevCorsEnabled bool
	DevCorsHost    string
}

type API struct {
	config  Config
	engine  *gin.Engine
	db      *sql.DB
	metrics *metrics.Metrics
}

func New(db *sql.DB, metrics *metrics.Metrics, config Config) *API {
	return &API{
		config:  config,
		db:      db,
		metrics: metrics,
	}
}

func (a *API) Start() {
	a.engine = gin.Default()

	if a.config.DevCorsEnabled {
		a.engine.Use(cors.New(cors.Config{
			AllowOrigins:     []string{a.config.DevCorsHost},
			AllowMethods:     []string{"PUT", "PATCH", "GET", "POST"},
			AllowHeaders:     []string{"Origin"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
		}))
	}

	a.setRoutes()

	err := a.engine.Run(":" + a.config.Port)
	if err != nil {
		log.Fatal(err)
	}
}

func (a *API) Close() {
}
