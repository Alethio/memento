package gui

import "github.com/gin-gonic/gin"

func (g *GUI) GUIResetHandler(c *gin.Context) {
	g.sendGUIResponse(c, "reset", gin.H{})
}

func (g *GUI) GUIResetPostHandler(c *gin.Context) {
	var errors []string
	var success []string

	defer func() {
		g.sendGUIResponse(c, "reset", gin.H{
			"errors":  errors,
			"success": success,
		})
	}()

	g.core.Pause()

	err := g.core.Reset()
	if err != nil {
		errors = append(errors, err.Error())
	} else {
		success = append(success, "The database has been reset.")
	}
}
