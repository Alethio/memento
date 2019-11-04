package dashboard

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (d *Dashboard) QueueHandler(c *gin.Context) {
	d.sendResponse(c, "queue", gin.H{})
}

func (d *Dashboard) QueuePostHandler(c *gin.Context) {
	var errors []string
	var success []string

	defer func() {
		d.sendResponse(c, "queue", gin.H{
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

		err = d.core.AddTodo(blockInt)
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
			err = d.core.AddTodo(i)
			if err != nil {
				errors = append(errors, fmt.Sprintf("Could not queue block: %s", err))
				return
			}
		}

		success = append(success, "Blocks successfully queued!")
	}
}
