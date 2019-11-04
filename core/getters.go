package core

import (
	"database/sql"

	"github.com/Alethio/memento/metrics"
)

func (c *Core) DB() *sql.DB {
	return c.db
}

func (c *Core) Metrics() *metrics.Provider {
	return c.metrics
}
