package api

import (
	"runtime"
	"strconv"

	"github.com/Alethio/memento/api/types"
	"github.com/Alethio/memento/data/storable"
)

func (a *API) getBlockTxs(number int64) ([]types.Tx, error) {
	var txs = make([]types.Tx, 0)

	rows, err := a.core.DB().Query(`select tx_index, tx_hash, value, "from", "to", msg_gas_limit, tx_gas_used, tx_gas_price from txs where included_in_block = $1 order by tx_index`, number)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	for rows.Next() {
		var (
			txIndex  int32
			txHash   string
			value    string
			from     storable.ByteArray
			to       storable.ByteArray
			gasLimit string
			gasUsed  string
			gasPrice string
		)

		err = rows.Scan(&txIndex, &txHash, &value, &from, &to, &gasLimit, &gasUsed, &gasPrice)
		if err != nil {
			log.Error(err)
			return nil, err
		}

		txs = append(txs, types.Tx{
			TxIndex:     &txIndex,
			TxHash:      &txHash,
			Value:       &value,
			From:        &from,
			To:          &to,
			MsgGasLimit: &gasLimit,
			TxGasUsed:   &gasUsed,
			TxGasPrice:  &gasPrice,
		})
	}

	return txs, nil
}

func (a *API) getDBEntries() (types.DBEntries, error) {
	var dbEntries types.DBEntries
	dbEntries = DBEntries

	err := a.core.DB().QueryRow(`
	select
	       (select count(*) from blocks)::text as blocks,
	       (select count(*) from txs)::text as txs,
	       (select count(*) from uncles)::text as uncles,
	       (select count(*) from log_entries)::text as log_entries
	`).Scan(&dbEntries.Blocks.Value, &dbEntries.Txs.Value, &dbEntries.Uncles.Value, &dbEntries.LogEntries.Value)
	if err != nil {
		log.Error(err)
		return dbEntries, err
	}

	return dbEntries, nil
}

func (a *API) getDBStats() (types.DBStats, error) {
	var dbStats types.DBStats
	dbStats = DBStats

	err := a.core.DB().QueryRow(`
		select pg_size_pretty(sum(table_size))   as table_size,
			   pg_size_pretty(sum(indexes_size)) as indexes_size,
			   pg_size_pretty(sum(total_size))   as total_size,
			   (select version_id from goose_db_version order by id desc limit 1) as migration_version,
		       coalesce((select number from blocks order by number desc limit 1)::text, 'null') as max_block
		from (
				 select table_name,
						pg_table_size(table_name)          as table_size,
						pg_indexes_size(table_name)        as indexes_size,
						pg_total_relation_size(table_name) as total_size
				 from (
						  select table_name::text
						  from information_schema.tables
						  where table_schema = 'public'
					  ) as all_tables
				 order by total_size desc
			 ) as pretty_sizes
     	`).Scan(&dbStats.DataSize.Value, &dbStats.IndexesSize.Value, &dbStats.TotalSize.Value, &dbStats.MigrationsVersion.Value, &dbStats.MaxBlock.Value)
	if err != nil {
		log.Error(err)
		return dbStats, err
	}

	return dbStats, nil
}

func (a *API) getProcStats() types.ProcStats {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	var procStats types.ProcStats
	procStats = ProcStats

	procStats.MemoryUsage.Value = strconv.FormatUint(bToMb(m.Sys), 10) + "MB"
	procStats.TodoLength.Value = strconv.FormatInt(a.core.Metrics().GetTodoLength(), 10)
	procStats.ReorgedBlocks.Value = strconv.FormatInt(a.core.Metrics().GetReorgedBlocks(), 10)
	procStats.InvalidBlocks.Value = strconv.FormatInt(a.core.Metrics().GetInvalidBlocks(), 10)

	return procStats
}

func (a *API) getTimingStats() types.TimingStats {
	var timingStats types.TimingStats
	timingStats = TimingStats

	timingStats.ProcessingTime.Value = a.core.Metrics().GetProcessingTime()
	timingStats.ScrapingTime.Value = a.core.Metrics().GetScrapingTime()
	timingStats.IndexingTime.Value = a.core.Metrics().GetIndexingTime()

	return timingStats
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
