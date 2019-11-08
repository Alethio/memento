package api

import (
	"encoding/hex"
	"fmt"

	"github.com/Alethio/memento/utils"
	"github.com/alethio/web3-go/ethrpc"
	"github.com/gin-gonic/gin"
)

func (a *API) AccountCodeHandler(c *gin.Context) {
	address := utils.CleanUpHex(c.Param("address"))
	if len(address) != 40 {
		BadRequest(c, fmt.Errorf("bad request: address is malformed"))
		return
	}

	eth, err := ethrpc.NewWithDefaults(a.config.EthClientURL)
	if err != nil {
		Error(c, err)
		return
	}

	code, err := eth.GetCode(fmt.Sprintf("0x%s", address))
	if err != nil {
		Error(c, err)
		return
	}

	OK(c, map[string]interface{}{
		"code": hex.EncodeToString(code),
	})
}

func (a *API) AccountBalanceHandler(c *gin.Context) {
	address := utils.CleanUpHex(c.Param("address"))
	if len(address) != 40 {
		BadRequest(c, fmt.Errorf("bad request: address is malformed"))
		return
	}

	eth, err := ethrpc.NewWithDefaults(a.config.EthClientURL)
	if err != nil {
		Error(c, err)
		return
	}

	balance, err := eth.GetBalanceAtBlock(fmt.Sprintf("0x%s", address), fmt.Sprintf("0x%x", a.core.Metrics().GetLatestBLock()))
	if err != nil {
		Error(c, err)
		return
	}

	OK(c, map[string]interface{}{
		"balance": balance.String(),
	})
}
