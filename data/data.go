package data

import (
	"database/sql"

	"github.com/Alethio/memento/data/storable"
	"github.com/sirupsen/logrus"

	"github.com/alethio/web3-go/types"
)

var log = logrus.WithField("module", "data")

type FullBlock struct {
	Block    types.Block
	Receipts []types.Receipt
	Uncles   []types.Block

	storables []Storable
}

// Storable
// role: a Storable serves as a means of transforming raw data and inserting it into the database
// input: raw Ethereum data + a database transaction
// output: processed/derived/enhanced data stored directly to the db
type Storable interface {
	ToDB(tx *sql.Tx) error
}

// RegisterStorables instantiates all the storables defined via code with the requested raw data
// Only the storables that are registered will be executed when the Store function is called
func (fb *FullBlock) RegisterStorables() {
	fb.storables = append(fb.storables, storable.NewStorableBlock(fb.Block))
	fb.storables = append(fb.storables, storable.NewStorableTxs(fb.Block, fb.Receipts))
	fb.storables = append(fb.storables, storable.NewStorableUncles(fb.Block, fb.Uncles))
	fb.storables = append(fb.storables, storable.NewStorableLogEntries(fb.Block, fb.Receipts))
	fb.storables = append(fb.storables, storable.NewStorableAccountTxs(fb.Block))
}

// Store will open a database transaction and execute all the registered Storables in the said transaction
func (fb *FullBlock) Store(db *sql.DB) error {
	exists, err := fb.checkBlockExists(db)
	if err != nil {
		return err
	}

	if exists {
		log.Info("block already exists in the database; skipping")
		return nil
	}

	reorged, err := fb.checkBlockReorged(db)
	if err != nil {
		return err
	}

	if reorged {
		number, err := fb.extractBlockNumber()
		if err != nil {
			return err
		}
		log.WithField("block", number).Warn("detected reorged block")
		_, err = db.Exec("select delete_block($1)", number)
		if err != nil {
			log.Error(err)
			return err
		}
		log.WithField("block", number).Info("removed old version from the db; will be replaced with new version")
	}

	tx, err := db.Begin()
	if err != nil {
		log.Error(err)
		return err
	}

	for _, s := range fb.storables {
		err = s.ToDB(tx)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}
