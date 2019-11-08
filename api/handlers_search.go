package api

import (
	"database/sql"
	"fmt"

	"github.com/Alethio/memento/utils"
	"github.com/gin-gonic/gin"
)

func (a *API) SearchHandler(c *gin.Context) {
	query := utils.CleanUpHex(c.Param("query"))

	if len(query) != 64 {
		BadRequest(c, fmt.Errorf("invalid request: invalid query string"))
		return
	}

	var count int
	err := a.core.DB().QueryRow(`select count(*) from txs where tx_hash = $1`, query).Scan(&count)
	if err != nil {
		Error(c, err)
		return
	}
	if count > 0 {
		OK(c, map[string]interface{}{
			"entity": "tx",
			"data":   nil,
		})
		return
	}

	var number int64
	err = a.core.DB().QueryRow(`select number from blocks where block_hash = $1`, query).Scan(&number)
	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}
	if err != sql.ErrNoRows {
		OK(c, map[string]interface{}{
			"entity": "block",
			"data": map[string]interface{}{
				"number": number,
			},
		})
		return
	}

	err = a.core.DB().QueryRow(`select count(*) from uncles where block_hash = $1`, query).Scan(&count)
	if err != nil {
		Error(c, err)
		return
	}
	if count > 0 {
		OK(c, map[string]interface{}{
			"entity": "uncle",
			"data":   nil,
		})
		return
	}

	OK(c, map[string]interface{}{})
}
