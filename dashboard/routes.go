package dashboard

import (
	"errors"
	"html/template"

	"github.com/gin-contrib/static"
)

func (d *Dashboard) setRoutes() {
	d.engine.SetFuncMap(template.FuncMap{
		"dict":  dict,
		"isMap": isMap,
		"plus":  plus,
		"times": times,
	})
	d.engine.LoadHTMLGlob("web/templates/**/*")
	d.engine.Use(static.Serve("/web/assets", static.LocalFile("web/assets", false)))
	d.engine.GET("/", d.IndexHandler)
	d.engine.GET("/queue", d.QueueHandler)
	d.engine.POST("/queue", d.QueuePostHandler)
	d.engine.GET("/pause", d.PauseHandler)
	d.engine.POST("/pause", d.PausePostHandler)
	d.engine.GET("/config", d.ConfigHandler)
	d.engine.POST("/config", d.ConfigPostHandler)
	d.engine.GET("/reset", d.ResetHandler)
	d.engine.POST("/reset", d.ResetPostHandler)
}

func dict(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, errors.New("invalid dict call")
	}

	dict := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, errors.New("dict keys must be strings")
		}
		dict[key] = values[i+1]
	}

	return dict, nil
}

func isMap(value interface{}) bool {
	_, ok := value.(map[string]interface{})
	return ok
}

func plus(x, y int) int {
	return x + y
}

func times(x, y int) int {
	return x * y
}
