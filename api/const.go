package api

import (
	"github.com/Alethio/memento/api/types"
)

const MaxBlocksInRange = 300

var DBEntries = types.DBEntries{
	Blocks:     types.Stat{Name: "Blocks", Icon: "cube-outline", Color: "green"},
	Txs:        types.Stat{Name: "Transactions", Icon: "swap-horizontal-circle-outline", Color: "green"},
	Uncles:     types.Stat{Name: "Uncles", Icon: "file-tree", Color: "green"},
	LogEntries: types.Stat{Name: "Log entries", Icon: "file-outline", Color: "green"},
}

var DBStats = types.DBStats{
	DataSize:          types.Stat{Name: "Data size", Icon: "file-document-box-outline", Color: "purple"},
	IndexesSize:       types.Stat{Name: "Indexes size", Icon: "table-search", Color: "purple"},
	TotalSize:         types.Stat{Name: "Total size", Icon: "database", Color: "purple"},
	MigrationsVersion: types.Stat{Name: "Migrations version", Icon: "source-branch", Color: "purple"},
	MaxBlock:          types.Stat{Name: "Highest processed block", Icon: "cube-send", Color: "blue"},
}

var ProcStats = types.ProcStats{
	MemoryUsage:   types.Stat{Name: "Current memory usage", Icon: "memory", Color: "orange"},
	TodoLength:    types.Stat{Name: "Tasks in todo", Icon: "clipboard-list-outline", Color: "blue"},
	ReorgedBlocks: types.Stat{Name: "Reorganized blocks", Icon: "link-variant-off", Color: "blue"},
	InvalidBlocks: types.Stat{Name: "Validation fails", Icon: "alert-circle-outline", Color: "blue"},
}

var TimingStats = types.TimingStats{
	ProcessingTime: types.Stat{Name: "Total time / block", Icon: "clock-outline", Color: "orange"},
	ScrapingTime:   types.Stat{Name: "Scraping time / block", Icon: "flip-horizontal", Color: "orange"},
	IndexingTime:   types.Stat{Name: "Indexing time / block", Icon: "card-search-outline", Color: "orange"},
}

var ViperIgnoredSettings = []string{"to", "from", "block", "version", "db.connection-string"}
