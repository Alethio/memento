package dashboard

import "github.com/gin-gonic/gin"

func (d *Dashboard) PauseHandler(c *gin.Context) {
	d.sendResponse(c, "pause", gin.H{})
}

func (d *Dashboard) PausePostHandler(c *gin.Context) {
	if d.core.IsPaused() {
		d.core.Resume()
	} else {
		d.core.Pause()
	}

	c.Redirect(302, "/pause")
}
