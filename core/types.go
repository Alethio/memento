package core

import (
	"git.aleth.io/alethio/memento/api"
	"git.aleth.io/alethio/memento/eth/bestblock"
	"git.aleth.io/alethio/memento/scraper"
	"git.aleth.io/alethio/memento/taskmanager"
)

type Features struct {
	Backfill    bool
	Lag         FeatureLag
	Automigrate bool
	Uncles      bool
}

type FeatureLag struct {
	Enabled bool
	Value   int64
}

type Config struct {
	BestBlockTracker         bestblock.Config
	TaskManager              taskmanager.Config
	Scraper                  scraper.Config
	PostgresConnectionString string
	API                      api.Config
	Features                 Features
}
