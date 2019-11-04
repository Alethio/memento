package dashboard

import "github.com/gin-gonic/gin"

func (d *Dashboard) ResetHandler(c *gin.Context) {
	d.sendResponse(c, "reset", gin.H{})
}

func (d *Dashboard) ResetPostHandler(c *gin.Context) {
	var errors []string
	var success []string

	defer func() {
		d.sendResponse(c, "reset", gin.H{
			"errors":  errors,
			"success": success,
		})
	}()

	d.core.Pause()

	err := d.core.Reset()
	if err != nil {
		errors = append(errors, err.Error())
	} else {
		success = append(success, "The database has been reset.")
	}
}
