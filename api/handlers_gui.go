package api

import (
	"fmt"
	"net/http"
	"strconv"

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

func (a *API) GUIQueueHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "queue", gin.H{
		"nav": types.Nav{
			Latest:  a.metrics.GetLatestBLock(),
			Version: viper.GetString("version"),
		},
	})
}

func (a *API) GUIQueuePostHandler(c *gin.Context) {
	var errors []string
	var success []string

	defer func() {
		c.HTML(http.StatusOK, "queue", gin.H{
			"nav": types.Nav{
				Latest:  a.metrics.GetLatestBLock(),
				Version: viper.GetString("version"),
			},
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

		err = a.taskmanager.Todo(blockInt)
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
			err = a.taskmanager.Todo(i)
			if err != nil {
				errors = append(errors, fmt.Sprintf("Could not queue block: %s", err))
				return
			}
		}

		success = append(success, "Blocks successfully queued!")
	}
}

func (a *API) GUIPauseHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "pause", gin.H{
		"nav": types.Nav{
			Latest:  a.metrics.GetLatestBLock(),
			Version: viper.GetString("version"),
		},
	})
}

func (a *API) GUIConfigHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "config", gin.H{
		"nav": types.Nav{
			Latest:  a.metrics.GetLatestBLock(),
			Version: viper.GetString("version"),
		},
	})
}

func (a *API) GUIResetHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "reset", gin.H{
		"nav": types.Nav{
			Latest:  a.metrics.GetLatestBLock(),
			Version: viper.GetString("version"),
		},
	})
}
