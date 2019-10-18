package data

import (
	"database/sql"
	"strconv"

	"github.com/Alethio/memento/data/storable"
)

// extractBlockNumber returns the block number as int64 by extracting it from the raw data
func (fb *FullBlock) extractBlockNumber() (int64, error) {
	number, err := strconv.ParseInt(fb.Block.Number, 0, 64)
	if err != nil {
		log.Error(err)
		return 0, err
	}

	return number, nil
}

// extractBlockHash returns the block hash as string by extracting it from the raw data
func (fb *FullBlock) extractBlockHash() string {
	return storable.Trim0x(fb.Block.Hash)
}

// checkBlockExists verifies if the current block matches any other block in the database by hash
func (fb *FullBlock) checkBlockExists(db *sql.DB) (bool, error) {
	hash := fb.extractBlockHash()

	var count int
	err := db.QueryRow(`select count(*) from blocks where block_hash = $1`, hash).Scan(&count)
	if err != nil {
		log.Error(err)
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

// checkBlockReorged verifies if the current block matches any block in the database on number
// this is meant to be used in order to detect if the database contains a blocks with the same number
// but different hash if the checkBlockExists function returns false
func (fb *FullBlock) checkBlockReorged(db *sql.DB) (bool, error) {
	number, err := fb.extractBlockNumber()
	if err != nil {
		return false, err
	}

	var count int
	err = db.QueryRow(`select count(*) from blocks where number = $1`, number).Scan(&count)
	if err != nil {
		log.Error(err)
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}
