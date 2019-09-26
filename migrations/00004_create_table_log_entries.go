package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upCreateTableLogEntries, downCreateTableLogEntries)
}

func upCreateTableLogEntries(tx *sql.Tx) error {
	_, err := tx.Exec(`
	create table log_entries
	(
		tx_hash                    text    not null,
		log_index                  integer not null,
		log_data                   bytea,
		logged_by                  text    not null,
		topic_0                    text,
		topic_1                    text,
		topic_2                    text,
		topic_3                    text,
		included_in_block          bigint  not null,
		created_at                 timestamp default now()
	);
	
	create index on log_entries (tx_hash, log_index);
	
	`)
	return err
}

func downCreateTableLogEntries(tx *sql.Tx) error {
	_, err := tx.Exec("drop table log_entries")
	return err
}
