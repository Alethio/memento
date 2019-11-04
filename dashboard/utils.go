package dashboard

import (
	"net/http"

	"github.com/Alethio/memento/dashboard/types"

	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"
)

func (d *Dashboard) sendResponse(c *gin.Context, template string, data gin.H) {
	c.HTML(http.StatusOK, template, mergeMaps(gin.H{
		"nav": types.Nav{
			Latest:  d.core.Metrics().GetLatestBLock(),
			Version: viper.GetString("version"),
			Paused:  d.core.IsPaused(),
		},
	}, data))
}

func mergeMaps(src, dst gin.H) gin.H {
	for k, v := range src {
		dst[k] = v
	}

	return dst
}
