package gui

import "github.com/gin-gonic/gin"

func (g *GUI) GUIPauseHandler(c *gin.Context) {
	g.sendGUIResponse(c, "pause", gin.H{})
}

func (g *GUI) GUIPausePostHandler(c *gin.Context) {
	if g.core.IsPaused() {
		g.core.Resume()
	} else {
		g.core.Pause()
	}

	c.Redirect(302, "/pause")
}
