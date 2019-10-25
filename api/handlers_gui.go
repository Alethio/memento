package api

import (
	"net/http"

	"github.com/Alethio/memento/api/types"

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

	c.HTML(http.StatusOK, "index", gin.H{
		"nav": types.Nav{
			Latest:  a.metrics.GetLatestBLock(),
			Version: viper.GetString("version"),
		},
		"dbEntries":   dbEntries,
		"dbStats":     dbStats,
		"procStats":   a.getProcStats(),
		"timingStats": a.getTimingStats(),
	})
}

func (a *API) GUISettingsHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "settings", gin.H{
		"nav": types.Nav{
			Latest:  a.metrics.GetLatestBLock(),
			Version: viper.GetString("version"),
		},
	})
}
