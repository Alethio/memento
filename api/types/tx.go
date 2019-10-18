package types

import "github.com/Alethio/memento/data/storable"

type Tx struct {
	TxHash              *string                      `json:"txHash,omitempty"`
	IncludedInBlock     *int64                       `json:"includedInBlock,omitempty"`
	TxIndex             *int32                       `json:"txIndex,omitempty"`
	From                *storable.ByteArray          `json:"from,omitempty"`
	To                  *storable.ByteArray          `json:"to,omitempty"`
	Value               *string                      `json:"value,omitempty"`
	TxNonce             *int64                       `json:"txNonce,omitempty"`
	MsgGasLimit         *string                      `json:"msgGasLimit,omitempty"`
	TxGasUsed           *string                      `json:"txGasUsed,omitempty"`
	TxGasPrice          *string                      `json:"txGasPrice,omitempty"`
	CumulativeGasUsed   *string                      `json:"cumulativeGasUsed,omitempty"`
	MsgPayload          *storable.ByteArray          `json:"msgPayload,omitempty"`
	MsgStatus           *string                      `json:"msgStatus,omitempty"`
	Creates             *storable.ByteArray          `json:"creates,omitempty"`
	TxLogsBloom         *storable.ByteArray          `json:"txLogsBloom,omitempty"`
	BlockCreationTime   *storable.DatetimeToJSONUnix `json:"blockCreationTime,omitempty"`
	LogEntriesTriggered *int32                       `json:"logEntriesTriggered,omitempty"`
}
