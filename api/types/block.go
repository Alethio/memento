package types

import "git.aleth.io/alethio/memento/data/storable"

type Block struct {
	Number               int64                       `json:"number"`
	BlockHash            string                      `json:"blockHash"`
	ParentBlockHash      string                      `json:"parentBlockHash"`
	BlockCreationTime    storable.DatetimeToJSONUnix `json:"blockCreationTime"`
	BlockGasLimit        string                      `json:"blockGasLimit"`
	BlockGasUsed         string                      `json:"blockGasUsed"`
	BlockDifficulty      string                      `json:"blockDifficulty"`
	TotalBlockDifficulty string                      `json:"totalBlockDifficulty"`
	BlockExtraData       storable.ByteArray          `json:"blockExtraData"`
	BlockMixHash         storable.ByteArray          `json:"blockMixHash"`
	BlockNonce           storable.ByteArray          `json:"blockNonce"`
	BlockSize            int64                       `json:"blockSize"`
	BlockLogsBloom       storable.ByteArray          `json:"blockLogsBloom"`
	IncludesUncle        storable.JSONStringArray    `json:"includesUncle"`
	HasBeneficiary       storable.ByteArray          `json:"hasBeneficiary"`
	HasReceiptsTrie      storable.ByteArray          `json:"hasReceiptsTrie"`
	HasTxTrie            storable.ByteArray          `json:"hasTxTrie"`
	Sha3Uncles           storable.ByteArray          `json:"sha3Uncles"`
	NumberOfUncles       int32                       `json:"numberOfUncles"`
	NumberOfTxs          int32                       `json:"numberOfTxs"`

	Txs []Tx `json:"txs"`
}
