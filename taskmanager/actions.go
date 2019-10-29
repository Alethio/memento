package taskmanager

import "github.com/go-redis/redis"

// Todo inserts a block into the redis sorted set used for queue management using a ZADD command
func (m *Manager) Todo(block int64) error {
	log.WithField("block", block).Trace("adding block to todo")
	return m.redis.ZAdd(m.config.TodoList, redis.Z{
		Score:  float64(block),
		Member: block,
	}).Err()
}

func (m *Manager) Reset() error {
	err := m.redis.Del(m.config.TodoList).Err()
	if err != nil {
		return err
	}

	// set the lastBlockAdded to -1 in order to backfill the whole chain after a reset if the backfill feature is enabled
	m.lastBlockAdded = -1

	return nil
}
