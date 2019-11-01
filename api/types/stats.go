package types

type DBEntries struct {
	Blocks, Txs, LogEntries, Uncles string
}

type DBStats struct {
	DataSize, IndexesSize, TotalSize, MigrationsVersion, MaxBlock string

	RawDataSize, RawIndexesSize int64
}

type ProcStats struct {
	ReorgedBlocks, InvalidBlocks, MemoryUsage, TodoLength, Version, PercentageDone string
}

type TimingStats struct {
	// human readable format
	ProcessingTime, ScrapingTime, IndexingTime string

	// value in ms
	RawProcessingTime, RawScrapingTime, RawIndexingTime int64
}

type Nav struct {
	Latest  int64
	Version string
	Paused  bool
}
