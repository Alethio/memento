package api

import (
	"database/sql"

	"github.com/Alethio/memento/api/types"
	"github.com/Alethio/memento/utils"
	"github.com/gin-gonic/gin"
)

func (a *API) UncleDetailsHandler(c *gin.Context) {
	blockHash := utils.CleanUpHex(c.Param("hash"))

	var uncle types.Uncle
	err := a.core.DB().QueryRow("select block_hash, included_in_block, number, block_creation_time, uncle_index, block_gas_limit, block_gas_used, has_beneficiary, block_difficulty, block_extra_data, block_mix_hash, block_nonce, sha3_uncles from uncles where block_hash = $1 limit 1", blockHash).Scan(&uncle.BlockHash, &uncle.IncludedInBlock, &uncle.Number, &uncle.BlockCreationTime, &uncle.UncleIndex, &uncle.BlockGasLimit, &uncle.BlockGasUsed, &uncle.HasBeneficiary, &uncle.BlockDifficulty, &uncle.BlockExtraData, &uncle.BlockMixHash, &uncle.BlockNonce, &uncle.Sha3Uncles)
	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}

	if err == sql.ErrNoRows {
		NotFound(c)
		return
	}

	OK(c, uncle)
}
