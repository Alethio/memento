package dashboard

import (
	"github.com/gin-gonic/gin"
)

func (d *Dashboard) IndexHandler(c *gin.Context) {
	var errors []string

	dbEntries, err := d.getDBEntries()
	if err != nil {
		errors = append(errors, err.Error())
	}

	dbStats, err := d.getDBStats()
	if err != nil {
		errors = append(errors, err.Error())
	}

	d.sendResponse(c, "index", gin.H{
		"dbEntries":   dbEntries,
		"dbStats":     dbStats,
		"procStats":   d.getProcStats(),
		"timingStats": d.getTimingStats(),
		"errors":      errors,
	})
}
