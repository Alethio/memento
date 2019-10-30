package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Alethio/memento/api/types"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func (a *API) GUIIndexHandler(c *gin.Context) {
	var errors []string

	dbEntries, err := a.getDBEntries()
	if err != nil {
		errors = append(errors, err.Error())
	}

	dbStats, err := a.getDBStats()
	if err != nil {
		errors = append(errors, err.Error())
	}

	a.sendGUIResponse(c, "index", gin.H{
		"dbEntries":   dbEntries,
		"dbStats":     dbStats,
		"procStats":   a.getProcStats(),
		"timingStats": a.getTimingStats(),
		"errors":      errors,
	})
}

func (a *API) GUIQueueHandler(c *gin.Context) {
	a.sendGUIResponse(c, "queue", gin.H{})
}

func (a *API) GUIQueuePostHandler(c *gin.Context) {
	var errors []string
	var success []string

	defer func() {
		a.sendGUIResponse(c, "queue", gin.H{
			"errors":  errors,
			"success": success,
		})
	}()

	if c.PostForm("type") == "single" {
		block := c.PostForm("block")
		blockInt, err := strconv.ParseInt(block, 10, 64)
		if err != nil {
			errors = append(errors, "Block number must be numeric!")
			return
		}

		err = a.core.AddTodo(blockInt)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Could not queue block: %s", err))
			return
		}

		success = append(success, "Block successfully queued!")
	} else {
		start := c.PostForm("start")
		startInt, err := strconv.ParseInt(start, 10, 64)
		if err != nil {
			errors = append(errors, "Start block must be numeric!")
			return
		}

		end := c.PostForm("end")
		endInt, err := strconv.ParseInt(end, 10, 64)
		if err != nil {
			errors = append(errors, "End block must be numeric!")
			return
		}

		for i := startInt; i <= endInt; i++ {
			err = a.core.AddTodo(i)
			if err != nil {
				errors = append(errors, fmt.Sprintf("Could not queue block: %s", err))
				return
			}
		}

		success = append(success, "Blocks successfully queued!")
	}
}

func (a *API) GUIPauseHandler(c *gin.Context) {
	a.sendGUIResponse(c, "pause", gin.H{})
}

func (a *API) GUIPausePostHandler(c *gin.Context) {
	if a.core.IsPaused() {
		a.core.Resume()
	} else {
		a.core.Pause()
	}

	c.Redirect(302, "/pause")
}

func (a *API) GUIConfigHandler(c *gin.Context) {
	var ignore = []string{"to", "from", "block", "version"}

	settings := viper.AllSettings()
	for _, v := range ignore {
		delete(settings, v)
	}

	a.sendGUIResponse(c, "config", gin.H{
		"settings": settings,
	})
}

func (a *API) GUIResetHandler(c *gin.Context) {
	a.sendGUIResponse(c, "reset", gin.H{})
}

func (a *API) GUIResetPostHandler(c *gin.Context) {
	var errors []string
	var success []string

	defer func() {
		a.sendGUIResponse(c, "reset", gin.H{
			"errors":  errors,
			"success": success,
		})
	}()

	a.core.Pause()

	err := a.core.Reset()
	if err != nil {
		errors = append(errors, err.Error())
	} else {
		success = append(success, "The database has been reset.")
	}
}

func (a *API) sendGUIResponse(c *gin.Context, template string, data gin.H) {
	c.HTML(http.StatusOK, template, mergeMaps(gin.H{
		"nav": types.Nav{
			Latest:  a.core.Metrics().GetLatestBLock(),
			Version: viper.GetString("version"),
			Paused:  a.core.IsPaused(),
		},
	}, data))
}

func mergeMaps(src, dst gin.H) gin.H {
	for k, v := range src {
		dst[k] = v
	}

	return dst
}
