package api

import (
	"errors"
	"html/template"

	"github.com/gin-contrib/static"
)

func (a *API) setRoutes() {
	explorer := a.engine.Group("/explorer")
	explorer.GET("/block/:block", a.BlockHandler)
	explorer.GET("/block-range/:start/:end", a.BlockRangeHandler)
	explorer.GET("/uncle/:hash", a.UncleDetailsHandler)
	explorer.GET("/tx/:txHash", a.TxDetailsHandler)
	explorer.GET("/tx/:txHash/log-entries", a.TxLogEntriesHandler)
	explorer.GET("/account/:address/txs", a.AccountTxsHandler)

	a.engine.SetFuncMap(template.FuncMap{
		"dict": dict,
	})
	a.engine.LoadHTMLGlob("web/templates/**/*")
	a.engine.Use(static.Serve("/web/assets", static.LocalFile("web/assets", false)))
	a.engine.GET("/", a.GUIIndexHandler)
	a.engine.GET("/queue", a.GUIQueueHandler)
	a.engine.GET("/pause", a.GUIPauseHandler)
	a.engine.GET("/config", a.GUIConfigHandler)
	a.engine.GET("/reset", a.GUIResetHandler)
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
