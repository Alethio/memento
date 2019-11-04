package gui

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func (g *GUI) GUIConfigHandler(c *gin.Context) {
	if viper.ConfigFileUsed() == "" {
		g.sendGUIResponse(c, "config", gin.H{
			"settings": map[string]interface{}{},
			"errors":   []string{"Memento did not start using a config file."},
		})

		return
	}

	g.sendGUIResponse(c, "config", gin.H{
		"settings": getSettings(),
	})
}

func (g *GUI) GUIConfigPostHandler(c *gin.Context) {
	if viper.ConfigFileUsed() == "" {
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

	disposableViper := viper.New()
	for k, v := range data {
		viper.Set(k, v)
		disposableViper.Set(k, v)
	}

	err := disposableViper.WriteConfigAs(viper.ConfigFileUsed())
	if err != nil {
		g.sendGUIResponse(c, "config", gin.H{
			"settings": getSettings(),
			"errors":   []string{err.Error()},
		})
		return
	}

	go g.core.ExitDelayed()

	g.sendGUIResponse(c, "config", gin.H{
		"settings": getSettings(),
		"success":  []string{"Config updated successfully. Application will be closed in 2 seconds to apply the new settings."},
	})
}
