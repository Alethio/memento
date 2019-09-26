package storable

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/lib/pq"

	"github.com/alethio/web3-go/types"
)

type AccountTxsGroup struct {
	RawBlock types.Block

	blockNumber int64

	accountTxs []*AccountTx
}

// AccountTx is a type of entity that represents a transaction between two accounts
// For each transaction, 2 AccountTx entities will be created, one for each direction of the transaction
// This helps with querying an account's transactions history in a specific order (e.g. chronological) and paginated
type AccountTx struct {
	Address         string
	Counterparty    string
	TxHash          string
	Out             bool
	IncludedInBlock int64
	TxIndex         int64
}

func NewStorableAccountTxs(block types.Block) *AccountTxsGroup {
	return &AccountTxsGroup{RawBlock: block}
}

func (atg *AccountTxsGroup) ToDB(tx *sql.Tx) error {
	if len(atg.RawBlock.Transactions) == 0 {
		return nil
	}

	log.Trace("storing account transactions")
	start := time.Now()
	defer func() {
		log.WithField("duration", time.Since(start)).WithField("count", len(atg.accountTxs)).Debug("done storing account transactions")
	}()

	err := atg.enhance()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(pq.CopyIn("account_txs", "address", "counterparty", "tx_hash", "out", "included_in_block", "tx_index"))
	if err != nil {
		return err
	}

	for _, at := range atg.accountTxs {
		_, err = stmt.Exec(at.Address, at.Counterparty, at.TxHash, at.Out, at.IncludedInBlock, at.TxIndex)
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

// enhance processes all the transactions in the raw block and generates a list containing the AccountTx entities
// corresponding to each of the block's transactions
func (atg *AccountTxsGroup) enhance() error {
	number, err := strconv.ParseInt(atg.RawBlock.Number, 0, 64)
	if err != nil {
		log.Error(err)
		return err
	}
	atg.blockNumber = number

	for _, tx := range atg.RawBlock.Transactions {
		txIndex, err := strconv.ParseInt(tx.TransactionIndex, 0, 64)
		if err != nil {
			log.Error(err)
			return err
		}

		storableAccountTxOut, err := atg.buildStorableAccountTx(tx, txIndex, true)
		if err != nil {
			return err
		}
		atg.accountTxs = append(atg.accountTxs, storableAccountTxOut)

		storableAccountTxIn, err := atg.buildStorableAccountTx(tx, txIndex, false)
		if err != nil {
			return err
		}
		atg.accountTxs = append(atg.accountTxs, storableAccountTxIn)
	}

	return nil
}

func (atg *AccountTxsGroup) buildStorableAccountTx(tx types.Transaction, txIndex int64, out bool) (*AccountTx, error) {
	at := &AccountTx{}
	at.IncludedInBlock = atg.blockNumber
	at.TxIndex = txIndex
	at.TxHash = Trim0x(tx.Hash)
	at.Out = out

	if out {
		at.Address = Trim0x(tx.From)
		at.Counterparty = Trim0x(tx.To)
	} else {
		at.Address = Trim0x(tx.To)
		at.Counterparty = Trim0x(tx.From)
	}

	return at, nil
}
