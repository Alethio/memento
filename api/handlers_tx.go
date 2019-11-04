package api

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/Alethio/memento/api/types"
	"github.com/Alethio/memento/data/storable"
	"github.com/Alethio/memento/utils"
	"github.com/gin-gonic/gin"
)

func (a *API) TxDetailsHandler(c *gin.Context) {
	searchHash := utils.CleanUpHex(c.Param("txHash"))

	var (
		txHash              string
		includedInBlock     int64
		txIndex             int32
		from                storable.ByteArray
		to                  storable.ByteArray
		value               string
		txNonce             int64
		msgGasLimit         string
		txGasUsed           string
		txGasPrice          string
		cumulativeGasUsed   string
		msgPayload          storable.ByteArray
		msgStatus           string
		creates             storable.ByteArray
		txLogsBloom         storable.ByteArray
		blockCreationTime   storable.DatetimeToJSONUnix
		logEntriesTriggered int32
	)
	err := a.core.DB().QueryRow(`select tx_hash, included_in_block, tx_index, "from", "to", value, tx_nonce, msg_gas_limit, tx_gas_used, tx_gas_price, cumulative_gas_used, msg_payload, msg_status, creates, tx_logs_bloom, block_creation_time, log_entries_triggered from txs where tx_hash = $1 limit 1`, searchHash).Scan(&txHash, &includedInBlock, &txIndex, &from, &to, &value, &txNonce, &msgGasLimit, &txGasUsed, &txGasPrice, &cumulativeGasUsed, &msgPayload, &msgStatus, &creates, &txLogsBloom, &blockCreationTime, &logEntriesTriggered)
	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}

	if err == sql.ErrNoRows {
		NotFound(c)
		return
	}

	OK(c, types.Tx{
		TxHash:              &txHash,
		IncludedInBlock:     &includedInBlock,
		TxIndex:             &txIndex,
		From:                &from,
		To:                  &to,
		Value:               &value,
		TxNonce:             &txNonce,
		MsgGasLimit:         &msgGasLimit,
		TxGasUsed:           &txGasUsed,
		TxGasPrice:          &txGasPrice,
		CumulativeGasUsed:   &cumulativeGasUsed,
		MsgPayload:          &msgPayload,
		MsgStatus:           &msgStatus,
		Creates:             &creates,
		TxLogsBloom:         &txLogsBloom,
		BlockCreationTime:   &blockCreationTime,
		LogEntriesTriggered: &logEntriesTriggered,
	})
}

func (a *API) TxLogEntriesHandler(c *gin.Context) {
	txHash := utils.CleanUpHex(c.Param("txHash"))

	rows, err := a.core.DB().Query(`select tx_hash, log_index, log_data, logged_by, topic_0, topic_1, topic_2, topic_3 from log_entries where tx_hash = $1 order by log_index`, txHash)
	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}

	var logEntries []types.LogEntry
	for rows.Next() {
		var (
			le                             types.LogEntry
			topic0, topic1, topic2, topic3 string
		)

		err := rows.Scan(&le.TxHash, &le.LogIndex, &le.LogData, &le.LoggedBy, &topic0, &topic1, &topic2, &topic3)
		if err != nil {
			Error(c, err)
			return
		}

		// we can do this since it is not possible, for example, to have an empty topic2 and a non-empty topic3
		le.HasLogTopics = utils.AppendNotEmpty(le.HasLogTopics, topic0)
		le.HasLogTopics = utils.AppendNotEmpty(le.HasLogTopics, topic1)
		le.HasLogTopics = utils.AppendNotEmpty(le.HasLogTopics, topic2)
		le.HasLogTopics = utils.AppendNotEmpty(le.HasLogTopics, topic3)

		le.EventDecoded = make(map[string]interface{})
		le.EventDecoded["topic0"] = fmt.Sprintf("0x%s", topic0)
		le.EventDecoded["event"] = ""

		var inputs []map[string]interface{}
		for i := 1; i <= 3; i++ {
			if len(le.HasLogTopics) > i {
				inputs = append(inputs, map[string]interface{}{
					"name":    fmt.Sprintf("topic%d", i),
					"type":    "raw",
					"indexed": true,
					"value":   le.HasLogTopics[i],
				})
			}
		}

		data := le.LogData.String()
		counter := 0
		for i := 0; i < len(data); i += 64 {
			inputs = append(inputs, map[string]interface{}{
				"name":  fmt.Sprintf("data%d", counter),
				"type":  "raw",
				"value": data[i : i+64],
			})
			counter++
		}
		le.EventDecoded["inputs"] = inputs

		logEntries = append(logEntries, le)
	}

	if len(logEntries) == 0 {
		NotFound(c)
		return
	}

	OK(c, logEntries)
}

func (a *API) AccountTxsHandler(c *gin.Context) {
	accountAddress, err := utils.ValidateAccount(c.Param("address"))
	if err != nil {
		BadRequest(c, err)
		return
	}

	limit := c.DefaultQuery("limit", "50")
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		BadRequest(c, fmt.Errorf("limit value should by numeric"))
		return
	}

	var includedInBlock, txIndex *int64

	index := c.Query("txIndex")
	if index != "" {
		indexInt, err := strconv.ParseInt(index, 10, 64)
		if err != nil {
			BadRequest(c, fmt.Errorf("txIndex must be a positive integer"))
			return
		}

		txIndex = &indexInt
	}

	block := c.Query("includedInBlock")
	if block != "" {
		blockInt, err := strconv.ParseInt(block, 10, 64)
		if err != nil {
			BadRequest(c, fmt.Errorf("includedInBlock must be a positive integer"))
			return
		}

		includedInBlock = &blockInt
	}

	var filters string
	var params []interface{}
	params = append(params, accountAddress, limitInt)

	if includedInBlock != nil {
		if txIndex == nil {
			zero := int64(0)
			txIndex = &zero
		}

		filters = "and (t1.included_in_block < $3 or (t1.included_in_block = $3 and t1.tx_index < $4))"
		params = append(params, includedInBlock, txIndex)
	}

	query := `select t2.tx_hash, t2.tx_index, "from", "to", value, block_creation_time, t2.included_in_block, tx_gas_used, tx_gas_price
				from account_txs as t1
				left join txs as t2 on (t2.tx_hash = t1.tx_hash)
				where t1.address = $1 %s
				order by t1.included_in_block desc, t1.tx_index desc limit $2`

	query = fmt.Sprintf(query, filters)

	rows, err := a.core.DB().Query(query, params...)
	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}

	var txs = make([]types.Tx, 0)
	for rows.Next() {
		var (
			txHash            string
			txIndex           int32
			from              storable.ByteArray
			to                storable.ByteArray
			value             string
			blockCreationTime storable.DatetimeToJSONUnix
			includedInBlock   int64
			txGasUsed         string
			txGasPrice        string
		)

		err := rows.Scan(&txHash, &txIndex, &from, &to, &value, &blockCreationTime, &includedInBlock, &txGasUsed, &txGasPrice)
		if err != nil {
			Error(c, err)
			return
		}

		txs = append(txs, types.Tx{
			TxHash:            &txHash,
			TxIndex:           &txIndex,
			From:              &from,
			To:                &to,
			Value:             &value,
			BlockCreationTime: &blockCreationTime,
			IncludedInBlock:   &includedInBlock,
			TxGasUsed:         &txGasUsed,
			TxGasPrice:        &txGasPrice,
		})
	}

	OK(c, txs)
}
