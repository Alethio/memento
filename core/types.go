package core

import (
	"github.com/Alethio/memento/api"
	"github.com/Alethio/memento/eth/bestblock"
	"github.com/Alethio/memento/scraper"
	"github.com/Alethio/memento/taskmanager"
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
