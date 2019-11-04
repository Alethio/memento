package core

import (
	"os"
	"time"
)

func (c *Core) AddTodo(block int64) error {
	return c.taskmanager.Todo(block)
}

func (c *Core) Pause() {
	c.taskmanager.Pause()
}

func (c *Core) Resume() {
	c.taskmanager.Resume()
}

func (c *Core) IsPaused() bool {
	return c.taskmanager.IsPaused()
}

func (c *Core) Reset() error {
	err := c.taskmanager.Reset()
	if err != nil {
		log.Error(err)
		return err
	}

	tx, err := c.db.Begin()
	if err != nil {
		log.Error(err)
		return err
	}

	_, err = tx.Exec(`
		truncate table blocks restart identity;
		truncate table uncles restart identity;
		truncate table txs restart identity;
		truncate table log_entries restart identity;
		truncate table account_txs restart identity;
		`)
	if err != nil {
		log.Error(err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Error(err)
		return err
	}

	c.metrics.Reset()
	c.metrics.RecordLatestBlock(c.bbtracker.BestBlock())

	return nil
}

func (c *Core) ExitDelayed() {
	c.Pause()
	time.Sleep(2 * time.Second)
	c.Close()
	os.Exit(0)
}
