package types

type Stat struct {
	Name  string
	Value string
	Icon  string
	Color string
}

type DBEntries struct {
	Blocks, Txs, LogEntries, Uncles Stat
}

type DBStats struct {
	DataSize, IndexesSize, TotalSize, MigrationsVersion, MaxBlock Stat
}

type ProcStats struct {
	ReorgedBlocks, InvalidBlocks, MemoryUsage, TodoLength, Version Stat
}

type TimingStats struct {
	ProcessingTime, ScrapingTime, IndexingTime Stat
}

type Nav struct {
	Latest  int64
	Version string
	Paused  bool
}
