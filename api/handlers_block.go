package api

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/Alethio/memento/api/types"

	"github.com/Alethio/memento/data/storable"

	"github.com/gin-gonic/gin"
)

func (a *API) BlockHandler(c *gin.Context) {
	var (
		blockNumber int64
		err         error
	)

	if c.Param("block") == "latest" {
		err := a.db.QueryRow("select number from blocks order by number desc limit 1").Scan(&blockNumber)
		if err != nil {
			Error(c, err)
			return
		}
	} else {
		blockNumber, err = strconv.ParseInt(c.Param("block"), 10, 64)
		if err != nil {
			BadRequest(c, fmt.Errorf("invalid request: block number must be numeric"))
			return
		}
	}

	var block types.Block
	err = a.db.QueryRow("select number, block_hash, parent_block_hash, block_creation_time, block_gas_limit, block_gas_used, block_difficulty, total_block_difficulty, block_extra_data, block_mix_hash, block_nonce, block_size, block_logs_bloom, includes_uncle, has_beneficiary, has_receipts_trie, has_tx_trie, sha3_uncles, number_of_uncles, number_of_txs from blocks where number = $1 limit 1", blockNumber).Scan(
		&block.Number,
		&block.BlockHash,
		&block.ParentBlockHash,
		&block.BlockCreationTime,
		&block.BlockGasLimit,
		&block.BlockGasUsed,
		&block.BlockDifficulty,
		&block.TotalBlockDifficulty,
		&block.BlockExtraData,
		&block.BlockMixHash,
		&block.BlockNonce,
		&block.BlockSize,
		&block.BlockLogsBloom,
		&block.IncludesUncle,
		&block.HasBeneficiary,
		&block.HasReceiptsTrie,
		&block.HasTxTrie,
		&block.Sha3Uncles,
		&block.NumberOfUncles,
		&block.NumberOfTxs,
	)
	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}

	if err == sql.ErrNoRows {
		NotFound(c)
		return
	}

	if len(block.IncludesUncle) == 0 {
		block.IncludesUncle = nil
	}

	block.Txs, err = a.getBlockTxs(blockNumber)
	if err != nil {
		Error(c, err)
		return
	}

	OK(c, block)
}

func (a *API) BlockRangeHandler(c *gin.Context) {
	start, err1 := strconv.ParseUint(c.Param("start"), 10, 64)
	end, err2 := strconv.ParseUint(c.Param("end"), 10, 64)
	if err1 != nil || err2 != nil {
		BadRequest(c, fmt.Errorf("invalid request: block number must be numeric"))
		return
	}

	if end <= start || end-start > MaxBlocksInRange {
		BadRequest(c, fmt.Errorf("invalid request: block range not valid"))
		return
	}

	rows, err := a.db.Query("select number, block_creation_time, has_beneficiary, number_of_txs from blocks where number between $1 and $2 order by number desc", start, end)
	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}

	var blockList = make([]map[string]interface{}, 0)

	for rows.Next() {
		var (
			number         int64
			creationTime   storable.DatetimeToJSONUnix
			hasBeneficiary storable.ByteArray
			numberOfTxs    int32
		)

		err := rows.Scan(&number, &creationTime, &hasBeneficiary, &numberOfTxs)
		if err != nil {
			Error(c, err)
			return
		}

		block := make(map[string]interface{})
		block["number"] = number
		block["blockCreationTime"] = creationTime
		block["hasBeneficiary"] = hasBeneficiary
		block["numberOfTxs"] = numberOfTxs

		blockList = append(blockList, block)
	}

	if len(blockList) == 0 {
		NotFound(c)
		return
	}

	OK(c, blockList)
}
