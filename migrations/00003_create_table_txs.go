package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upCreateTableTxs, downCreateTableTxs)
}

func upCreateTableTxs(tx *sql.Tx) error {
	_, err := tx.Exec(`
	create table txs (
		tx_hash                        text        not null,
		included_in_block              bigint      not null,
		tx_index                       integer     not null,
		"from"                         bytea       not null,
		"to"                           bytea       not null,
		value                          numeric(28) not null,
		tx_nonce                       bigint      not null,
		msg_gas_limit                  numeric(28) not null,
		tx_gas_used                    numeric(28),
		tx_gas_price                   numeric(28) not null,
		cumulative_gas_used            numeric(28) not null,
		msg_payload                    bytea,
		msg_status                     text,
		creates                        bytea,
		tx_logs_bloom                  bytea,
		block_creation_time            timestamp with time zone,
		log_entries_triggered          integer   default 0,
		created_at                     timestamp default now()
	);
	
	create index on txs (tx_hash);
	create index on txs (included_in_block desc, tx_index desc);
	
	`)
	return err
}

func downCreateTableTxs(tx *sql.Tx) error {
	_, err := tx.Exec("drop table txs;")
	return err
}
