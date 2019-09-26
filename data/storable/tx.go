package storable

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/lib/pq"

	"github.com/alethio/web3-go/types"
)

type TxsGroup struct {
	RawBlock    types.Block
	RawReceipts []types.Receipt

	blockNumber       int64
	blockCreationTime DatetimeToJSONUnix
	txs               []*Tx
}

type Tx struct {
	TxHash              string
	IncludedInBlock     int64
	TxIndex             int32
	From                ByteArray
	To                  ByteArray
	Value               string
	TxNonce             int64
	MsgGasLimit         string
	TxGasUsed           string
	TxGasPrice          string
	CumulativeGasUsed   string
	MsgPayload          ByteArray
	MsgStatus           string
	Creates             ByteArray
	TxLogsBloom         ByteArray
	BlockCreationTime   DatetimeToJSONUnix
	LogEntriesTriggered int32
}

func NewStorableTxs(block types.Block, receipts []types.Receipt) *TxsGroup {
	return &TxsGroup{RawBlock: block, RawReceipts: receipts}
}

func (t *TxsGroup) ToDB(dbTx *sql.Tx) error {
	if len(t.RawBlock.Transactions) == 0 {
		return nil
	}

	log.Trace("storing txs")
	start := time.Now()
	defer func() {
		log.WithField("duration", time.Since(start)).WithField("count", len(t.txs)).Debug("done storing txs")
	}()

	err := t.enhance()
	if err != nil {
		return err
	}

	stmt, err := dbTx.Prepare(pq.CopyIn("txs", "tx_hash", "included_in_block", "tx_index", "from", "to", "value", "tx_nonce", "msg_gas_limit", "tx_gas_used", "tx_gas_price", "cumulative_gas_used", "msg_payload", "msg_status", "creates", "tx_logs_bloom", "block_creation_time", "log_entries_triggered"))
	if err != nil {
		return err
	}

	for _, tx := range t.txs {
		_, err = stmt.Exec(tx.TxHash, tx.IncludedInBlock, tx.TxIndex, tx.From, tx.To, tx.Value, tx.TxNonce, tx.MsgGasLimit, tx.TxGasUsed, tx.TxGasPrice, tx.CumulativeGasUsed, tx.MsgPayload, tx.MsgStatus, tx.Creates, tx.TxLogsBloom, tx.BlockCreationTime, tx.LogEntriesTriggered)
		if err != nil {
			return err
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	err = stmt.Close()
	if err != nil {
		return err
	}

	return nil
}

// enhance processes the raw block and raw receipts data generating a list containing all the txs in the block
// in a format that's ready for insertion into the database
func (t *TxsGroup) enhance() error {
	number, err := strconv.ParseInt(t.RawBlock.Number, 0, 64)
	if err != nil {
		log.Error(err)
		return err
	}
	t.blockNumber = number

	timestamp, err := strconv.ParseInt(t.RawBlock.Timestamp, 0, 64)
	if err != nil {
		log.Error(err)
		return err
	}
	t.blockCreationTime = DatetimeToJSONUnix(time.Unix(timestamp, 0))

	for index, tx := range t.RawBlock.Transactions {
		storableTx, err := t.buildStorableTx(tx, t.RawReceipts[index])
		if err != nil {
			return err
		}

		t.txs = append(t.txs, storableTx)
	}

	return nil
}

func (t *TxsGroup) buildStorableTx(tx types.Transaction, receipt types.Receipt) (*Tx, error) {
	sTx := &Tx{}
	sTx.IncludedInBlock = t.blockNumber
	sTx.BlockCreationTime = t.blockCreationTime

	sTx.TxHash = Trim0x(tx.Hash)
	sTx.From = ByteArray(Trim0x(tx.From))
	sTx.To = ByteArray(Trim0x(tx.To))
	if tx.To == "" {
		if contractAddress, ok := receipt.ContractAddress.(string); ok && contractAddress != "" {
			sTx.To = ByteArray(Trim0x(contractAddress))
			sTx.Creates = ByteArray(Trim0x(contractAddress))
		}
	}

	if tx.Creates != "" {
		sTx.Creates = ByteArray(Trim0x(tx.Creates))
	}

	sTx.MsgPayload = ByteArray(Trim0x(tx.Input))
	sTx.TxLogsBloom = ByteArray(Trim0x(receipt.LogsBloom))
	sTx.MsgStatus = receipt.Status

	// -- int
	txIndex, err := strconv.ParseInt(tx.TransactionIndex, 0, 32)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	sTx.TxIndex = int32(txIndex)

	txNonce, err := strconv.ParseInt(tx.Nonce, 0, 64)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	sTx.TxNonce = txNonce

	// -- bigint
	gasLimit, err := HexStrToBigIntStr(tx.Gas)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	sTx.MsgGasLimit = gasLimit

	value, err := HexStrToBigIntStr(tx.Value)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	sTx.Value = value

	gasUsed, err := HexStrToBigIntStr(receipt.GasUsed)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	sTx.TxGasUsed = gasUsed

	gasPrice, err := HexStrToBigIntStr(tx.GasPrice)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	sTx.TxGasPrice = gasPrice

	cumulativeGasUsed, err := HexStrToBigIntStr(receipt.CumulativeGasUsed)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	sTx.CumulativeGasUsed = cumulativeGasUsed

	// -- computed
	sTx.LogEntriesTriggered = int32(len(receipt.Logs))

	return sTx, nil
}
