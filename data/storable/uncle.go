package storable

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/lib/pq"

	"github.com/alethio/web3-go/types"
)

type UnclesGroup struct {
	RawBlock  types.Block
	RawUncles []types.Block

	blockNumber int64

	uncles []*Uncle
}

type Uncle struct {
	BlockHash         string
	IncludedInBlock   int64
	Number            int64
	BlockCreationTime DatetimeToJSONUnix
	UncleIndex        int32
	BlockGasLimit     string
	BlockGasUsed      string
	HasBeneficiary    ByteArray
	BlockDifficulty   string
	BlockExtraData    ByteArray
	BlockMixHash      ByteArray
	BlockNonce        ByteArray
	Sha3Uncles        ByteArray
}

func NewStorableUncles(block types.Block, uncles []types.Block) *UnclesGroup {
	return &UnclesGroup{RawBlock: block, RawUncles: uncles}
}

func (ug *UnclesGroup) ToDB(tx *sql.Tx) error {
	if len(ug.RawUncles) == 0 {
		return nil
	}

	log.Trace("storing uncles")
	start := time.Now()
	defer func() {
		log.WithField("duration", time.Since(start)).WithField("count", len(ug.uncles)).Debug("done storing uncles")
	}()

	err := ug.enhance()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(pq.CopyIn("uncles", "block_hash", "included_in_block", "number", "block_creation_time", "uncle_index", "block_gas_limit", "block_gas_used", "has_beneficiary", "block_difficulty", "block_extra_data", "block_mix_hash", "block_nonce", "sha3_uncles"))
	if err != nil {
		return err
	}

	for _, uncle := range ug.uncles {
		_, err = stmt.Exec(uncle.BlockHash, uncle.IncludedInBlock, uncle.Number, uncle.BlockCreationTime, uncle.UncleIndex, uncle.BlockGasLimit, uncle.BlockGasUsed, uncle.HasBeneficiary, uncle.BlockDifficulty, uncle.BlockExtraData, uncle.BlockMixHash, uncle.BlockNonce, uncle.Sha3Uncles)
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

func (ug *UnclesGroup) enhance() error {
	number, err := strconv.ParseInt(ug.RawBlock.Number, 0, 64)
	if err != nil {
		log.Error(err)
		return err
	}
	ug.blockNumber = number

	for index, uncle := range ug.RawUncles {
		storableUncle, err := ug.buildStorableUncle(uncle, int32(index))
		if err != nil {
			return err
		}

		ug.uncles = append(ug.uncles, storableUncle)
	}

	return nil
}

func (ug *UnclesGroup) buildStorableUncle(uncle types.Block, index int32) (*Uncle, error) {
	u := &Uncle{}
	u.IncludedInBlock = ug.blockNumber
	u.UncleIndex = index

	if uncle.Miner == "" {
		uncle.Miner = uncle.Author
	}

	// -- raw
	u.BlockHash = Trim0x(uncle.Hash)
	u.HasBeneficiary = ByteArray(Trim0x(uncle.Miner))
	u.BlockExtraData = ByteArray(Trim0x(uncle.ExtraData))
	u.BlockMixHash = ByteArray(Trim0x(uncle.MixHash))
	u.BlockNonce = ByteArray(Trim0x(uncle.Nonce))
	u.Sha3Uncles = ByteArray(Trim0x(uncle.Sha3Uncles))

	// -- int64
	number, err := strconv.ParseInt(uncle.Number, 0, 64)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	u.Number = number

	// -- hexes
	gasLimit, err := HexStrToBigIntStr(uncle.GasLimit)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	u.BlockGasLimit = gasLimit

	gasUsed, err := HexStrToBigIntStr(uncle.GasUsed)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	u.BlockGasUsed = gasUsed

	difficulty, err := HexStrToBigIntStr(uncle.Difficulty)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	u.BlockDifficulty = difficulty

	// -- timestamp
	timestamp, err := strconv.ParseInt(uncle.Timestamp, 0, 64)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	u.BlockCreationTime = DatetimeToJSONUnix(time.Unix(timestamp, 0))

	return u, nil
}
