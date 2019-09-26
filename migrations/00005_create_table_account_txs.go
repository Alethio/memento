package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upCreateTableAccountTxs, downCreateTableAccountTxs)
}

func upCreateTableAccountTxs(tx *sql.Tx) error {
	_, err := tx.Exec(`
	create table account_txs
	(
		address 			text not null,
		counterparty 		text not null,
		tx_hash 			text not null,
		out 				bool,
		included_in_block 	bigint not null,
		tx_index 			bigint not null
	);
	
	create index on account_txs (address, included_in_block desc, tx_index desc);
	create index on account_txs (included_in_block desc);
	
	`)
	return err
}

func downCreateTableAccountTxs(tx *sql.Tx) error {
	_, err := tx.Exec("drop table account_txs;")
	return err
}
