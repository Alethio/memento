package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upCreateTableBlocks, downCreateTableBlocks)
}

func upCreateTableBlocks(tx *sql.Tx) error {
	_, err := tx.Exec(`
	create table blocks (
		number                   bigint      not null,
		block_hash               text        not null,
		parent_block_hash        text        not null,
		block_creation_time      timestamp with time zone,
		block_gas_limit          numeric(78) not null,
		block_gas_used           numeric(78) not null,
		block_difficulty         numeric(78) not null,
		total_block_difficulty   numeric(78) not null,
		block_extra_data         bytea,
		block_mix_hash           bytea       not null,
		block_nonce              bytea       not null,
		block_size               bigint      not null,
		block_logs_bloom         bytea       not null,
		includes_uncle           json,
		has_beneficiary          bytea,
		has_receipts_trie        bytea,
		has_tx_trie              bytea,
		sha3_uncles              bytea,
		number_of_uncles         integer     default 0,
		number_of_txs            integer     default 0,
		created_at               timestamp   default now()
	);
	
	create index on blocks (block_hash);
	create index on blocks (number);
	
	`)
	return err
}

func downCreateTableBlocks(tx *sql.Tx) error {
	_, err := tx.Exec("drop table blocks;")
	return err
}
