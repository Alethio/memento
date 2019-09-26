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
