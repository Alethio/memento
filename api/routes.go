package api

import "github.com/gin-gonic/gin"

func (a *API) setRoutes() {
	http := a.engine.Group("/api")
	http.GET("/block/:block", a.BlockHandler)
	http.GET("/block-range/:start/:end", a.BlockRangeHandler)
	http.GET("/uncle/:hash", a.UncleDetailsHandler)
	http.GET("/tx/:txHash", a.TxDetailsHandler)
	http.GET("/tx/:txHash/log-entries", a.TxLogEntriesHandler)
	http.GET("/account/:address/txs", a.AccountTxsHandler)

	a.engine.GET("/", func(ctx *gin.Context) {
		ctx.String(200, "It works!")
	})
}
