package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(Up00006, Down00006)
}

func Up00006(tx *sql.Tx) error {
	_, err := tx.Exec(`
	create or replace function __delete_entity(in table_name varchar, in number bigint) returns void as
	$body$
	begin
		execute format('delete from %1$s where included_in_block = %2$s;', table_name, number);
	end;
	$body$ language 'plpgsql';

	create or replace function delete_block(in block_number bigint) returns void as
	$body$
	declare
		tables varchar[];
		tbl    varchar;
	begin
		tables := array [
			'uncles',
			'txs',
			'log_entries',
			'account_txs'
			];
	
		foreach tbl in array tables
			loop
				perform __delete_entity(tbl, block_number);
			end loop;

		delete from blocks where number = block_number;
	end;
	$body$ language 'plpgsql';
	`)
	return err
}

func Down00006(tx *sql.Tx) error {
	_, err := tx.Exec(`
	drop function __delete_entity;
	drop function delete_block;
	`)
	return err
}
