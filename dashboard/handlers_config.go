package dashboard

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func (d *Dashboard) ConfigHandler(c *gin.Context) {
	if !d.config.ConfigEnabled {
		d.sendResponse(c, "config", gin.H{
			"disabled": true,
		})

		return
	}

	if viper.ConfigFileUsed() == "" {
		d.sendResponse(c, "config", gin.H{
			"settings": map[string]interface{}{},
			"errors":   []string{"Memento did not start using a config file."},
		})

		return
	}

	d.sendResponse(c, "config", gin.H{
		"settings": getSettings(),
	})
}

func (d *Dashboard) ConfigPostHandler(c *gin.Context) {
	if viper.ConfigFileUsed() == "" || !d.config.ConfigEnabled {
		c.Redirect(302, "/config")

		return
	}

	var data = make(map[string]interface{})

	for _, k := range viper.AllKeys() {
		v, exists := c.GetPostForm(fmt.Sprintf(".%s", k))

		// booleans are treated as a toggle (checkbox behind the scenes) in the interface
		// we're doing this due to the checkboxes behavior = not sending data if they're unchecked
		if _, ok := viper.Get(k).(bool); ok {
			data[k] = exists && v == "on"
			continue
		}

		if !exists {
			continue
		}

		data[k] = v
	}

	for _, k := range ViperIgnoredSettings {
		delete(data, k)
	}

	if data["api.port"] == data["dashboard.port"] {
		d.sendResponse(c, "config", gin.H{
			"settings": getSettings(),
			"errors":   []string{"API port can't be the same with the Dashboard port!"},
		})
		return
	}

	disposableViper := viper.New()
	for k, v := range data {
		viper.Set(k, v)
		disposableViper.Set(k, v)
	}

	err := disposableViper.WriteConfigAs(viper.ConfigFileUsed())
	if err != nil {
		d.sendResponse(c, "config", gin.H{
			"settings": getSettings(),
			"errors":   []string{err.Error()},
		})
		return
	}

	go d.core.ExitDelayed()

	d.sendResponse(c, "config", gin.H{
		"settings": getSettings(),
		"success":  []string{"Config updated successfully. Application will be closed in 2 seconds to apply the new settings."},
	})
}
