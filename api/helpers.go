package api

import (
	"fmt"
	"runtime"
	"strconv"

	"github.com/spf13/viper"

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

	err := a.core.DB().QueryRow(`
	select
	       (select count(*) from blocks)::text as blocks,
	       (select count(*) from txs)::text as txs,
	       (select count(*) from uncles)::text as uncles,
	       (select count(*) from log_entries)::text as log_entries
	`).Scan(&dbEntries.Blocks, &dbEntries.Txs, &dbEntries.Uncles, &dbEntries.LogEntries)
	if err != nil {
		log.Error(err)
		return dbEntries, err
	}

	return dbEntries, nil
}

func (a *API) getDBStats() (types.DBStats, error) {
	var dbStats types.DBStats

	err := a.core.DB().QueryRow(`
		select pg_size_pretty(sum(table_size))   as table_size,
		       sum(table_size) 					 as raw_table_size,
			   pg_size_pretty(sum(indexes_size)) as indexes_size,
		       sum(indexes_size) 				 as raw_indexes_size,
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
     	`).Scan(&dbStats.DataSize, &dbStats.RawDataSize, &dbStats.IndexesSize, &dbStats.RawIndexesSize, &dbStats.TotalSize, &dbStats.MigrationsVersion, &dbStats.MaxBlock)
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

	procStats.MemoryUsage = strconv.FormatUint(bToMb(m.Sys), 10) + "MB"
	procStats.TodoLength = strconv.FormatInt(a.core.Metrics().GetTodoLength(), 10)
	procStats.ReorgedBlocks = strconv.FormatInt(a.core.Metrics().GetReorgedBlocks(), 10)
	procStats.InvalidBlocks = strconv.FormatInt(a.core.Metrics().GetInvalidBlocks(), 10)

	procStats.PercentageDone = fmt.Sprintf("%f", 1-float64(a.core.Metrics().GetTodoLength())/float64(a.core.Metrics().GetLatestBLock()))

	return procStats
}

func (a *API) getTimingStats() types.TimingStats {
	var timingStats types.TimingStats

	timingStats.ProcessingTime = a.core.Metrics().GetProcessingTime()
	timingStats.RawProcessingTime = a.core.Metrics().GetRawProcessingTime()

	timingStats.ScrapingTime = a.core.Metrics().GetScrapingTime()
	timingStats.RawScrapingTime = a.core.Metrics().GetRawScrapingTime()

	timingStats.IndexingTime = a.core.Metrics().GetIndexingTime()
	timingStats.RawIndexingTime = a.core.Metrics().GetRawIndexingTime()

	return timingStats
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func getSettings() map[string]interface{} {
	settings := viper.AllSettings()
	for _, v := range ViperIgnoredSettings {
		delete(settings, v)
	}

	delete(settings["db"].(map[string]interface{}), "connection-string")

	return settings
}
