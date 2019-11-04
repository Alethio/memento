package gui

import (
	"github.com/gin-gonic/gin"
)

func (g *GUI) IndexHandler(c *gin.Context) {
	var errors []string

	dbEntries, err := g.getDBEntries()
	if err != nil {
		errors = append(errors, err.Error())
	}

	dbStats, err := g.getDBStats()
	if err != nil {
		errors = append(errors, err.Error())
	}

	g.sendGUIResponse(c, "index", gin.H{
		"dbEntries":   dbEntries,
		"dbStats":     dbStats,
		"procStats":   g.getProcStats(),
		"timingStats": g.getTimingStats(),
		"errors":      errors,
	})
}
