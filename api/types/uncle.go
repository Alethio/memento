package types

import "git.aleth.io/alethio/memento/data/storable"

type Uncle struct {
	BlockHash         string                      `json:"blockHash"`
	IncludedInBlock   int64                       `json:"includedInBlock"`
	Number            int64                       `json:"number"`
	BlockCreationTime storable.DatetimeToJSONUnix `json:"blockCreationTime"`
	UncleIndex        int32                       `json:"uncleIndex"`
	BlockGasLimit     string                      `json:"blockGasLimit"`
	BlockGasUsed      string                      `json:"blockGasUsed"`
	HasBeneficiary    storable.ByteArray          `json:"hasBeneficiary"`
	BlockDifficulty   string                      `json:"blockDifficulty"`
	BlockExtraData    storable.ByteArray          `json:"blockExtraData"`
	BlockMixHash      storable.ByteArray          `json:"blockMixHash"`
	BlockNonce        storable.ByteArray          `json:"blockNonce"`
	Sha3Uncles        storable.ByteArray          `json:"sha3Uncles"`
}
