package gui

import (
	"errors"
	"html/template"

	"github.com/gin-contrib/static"
)

func (g *GUI) setRoutes() {
	g.engine.SetFuncMap(template.FuncMap{
		"dict":  dict,
		"isMap": isMap,
		"plus":  plus,
		"times": times,
	})
	g.engine.LoadHTMLGlob("web/templates/**/*")
	g.engine.Use(static.Serve("/web/assets", static.LocalFile("web/assets", false)))
	g.engine.GET("/", g.IndexHandler)
	g.engine.GET("/queue", g.QueueHandler)
	g.engine.POST("/queue", g.QueuePostHandler)
	g.engine.GET("/pause", g.GUIPauseHandler)
	g.engine.POST("/pause", g.GUIPausePostHandler)
	g.engine.GET("/config", g.GUIConfigHandler)
	g.engine.POST("/config", g.GUIConfigPostHandler)
	g.engine.GET("/reset", g.GUIResetHandler)
	g.engine.POST("/reset", g.GUIResetPostHandler)
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
