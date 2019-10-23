package api

import (
	"net/http"

	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"
)

func (a *API) GUIIndexHandler(c *gin.Context) {
	dbEntries, err := a.getDBEntries()
	if err != nil {
		//
	}

	dbStats, err := a.getDBStats()
	if err != nil {
		//
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"latest":      a.metrics.GetLatestBLock(),
		"version":     viper.GetString("version"),
		"dbEntries":   dbEntries,
		"dbStats":     dbStats,
		"procStats":   a.getProcStats(),
		"timingStats": a.getTimingStats(),
	})
}
