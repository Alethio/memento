package dashboard

import (
	"fmt"
	"runtime"
	"strconv"

	"github.com/spf13/viper"

	"github.com/Alethio/memento/dashboard/types"
)

func (d *Dashboard) getDBEntries() (types.DBEntries, error) {
	var dbEntries types.DBEntries

	err := d.core.DB().QueryRow(`
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

func (d *Dashboard) getDBStats() (types.DBStats, error) {
	var dbStats types.DBStats

	err := d.core.DB().QueryRow(`
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

func (d *Dashboard) getProcStats() types.ProcStats {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	var procStats types.ProcStats

	procStats.MemoryUsage = strconv.FormatUint(bToMb(m.Sys), 10) + "MB"
	procStats.TodoLength = strconv.FormatInt(d.core.Metrics().GetTodoLength(), 10)
	procStats.ReorgedBlocks = strconv.FormatInt(d.core.Metrics().GetReorgedBlocks(), 10)
	procStats.InvalidBlocks = strconv.FormatInt(d.core.Metrics().GetInvalidBlocks(), 10)

	procStats.PercentageDone = fmt.Sprintf("%f", 1-float64(d.core.Metrics().GetTodoLength())/float64(d.core.Metrics().GetLatestBLock()))

	return procStats
}

func (d *Dashboard) getTimingStats() types.TimingStats {
	var timingStats types.TimingStats

	timingStats.ProcessingTime = d.core.Metrics().GetProcessingTime()
	timingStats.RawProcessingTime = d.core.Metrics().GetRawProcessingTime()

	timingStats.ScrapingTime = d.core.Metrics().GetScrapingTime()
	timingStats.RawScrapingTime = d.core.Metrics().GetRawScrapingTime()

	timingStats.IndexingTime = d.core.Metrics().GetIndexingTime()
	timingStats.RawIndexingTime = d.core.Metrics().GetRawIndexingTime()

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
