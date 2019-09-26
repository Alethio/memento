package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upCreateTableUncles, downCreateTableUncles)
}

func upCreateTableUncles(tx *sql.Tx) error {
	_, err := tx.Exec(`
	create table uncles (
		block_hash               text                     not null,
		included_in_block        bigint                   not null,
		number                   bigint                   not null,
		block_creation_time      timestamp with time zone not null,
		uncle_index              integer                  not null,
		block_gas_limit          numeric(78)              not null,
		block_gas_used           numeric(78)              not null,
		has_beneficiary          bytea                    not null,
		block_difficulty         numeric(78)              not null,
		block_extra_data         bytea                    not null,
		block_mix_hash           bytea                    not null,
		block_nonce              bytea                    not null,
		sha3_uncles              bytea                    not null,
		created_at               timestamp default now()
	);
	
	create index on uncles (block_hash);

	`)
	return err
}

func downCreateTableUncles(tx *sql.Tx) error {
	_, err := tx.Exec("drop table uncles;")
	return err
}
