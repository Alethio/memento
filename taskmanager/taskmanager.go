package taskmanager

import (
	"strconv"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/Alethio/memento/eth/bestblock"
	"github.com/go-redis/redis"
)

var log = logrus.WithField("module", "taskmanager")

type Config struct {
	RedisServer   string
	RedisPassword string
	TodoList      string
}

type Manager struct {
	config Config
	lag    int64

	tracker *bestblock.Tracker
	redis   *redis.Client

	closed   bool
	stopChan chan bool
}

// New instantiates a new task manager and also takes care of the redis connection management
// it subscribes to the best block tracker for new blocks which it'll add to the redis queue automatically
func New(tracker *bestblock.Tracker, lag int64, config Config) (*Manager, error) {
	m := &Manager{
		config:   config,
		lag:      lag,
		tracker:  tracker,
		closed:   false,
		stopChan: make(chan bool),
	}

	log.Info("setting up redis connection")
	m.redis = redis.NewClient(&redis.Options{
		Addr:        config.RedisServer,
		Password:    config.RedisPassword,
		DB:          0,
		ReadTimeout: time.Second * 1,
	})

	err := m.redis.Ping().Err()
	if err != nil {
		log.Error(err)
		return nil, err
	}

	log.Info("connected to redis successfully")

	go m.watchNewBlocks()

	return m, nil
}

func (m *Manager) Close() error {
	m.closed = true
	m.stopChan <- true
	return m.redis.Close()
}

// FeedToChan continuously executes blocking pops from the redis queue and sends the resulting task on the provided channel
// Highest blocks have priority
func (m *Manager) FeedToChan(c chan int64) {
	log.WithField("list", m.config.TodoList).Trace("feeding tasks from redis")
	for {
		if m.closed {
			if len(m.stopChan) > 0 {
				<-m.stopChan
			}
			return
		}

		doneChan := make(chan bool)
		var taskInt int64
		go func() {
			taskResult, err := m.redis.BZPopMax(0, m.config.TodoList).Result()
			if err != nil && m.closed {
				return
			}
			if err != nil {
				log.Error("getting task from redis returned error:", err)
				doneChan <- false
			}

			taskInt, err = strconv.ParseInt(taskResult.Member.(string), 10, 64)
			if err != nil {
				log.Error(err)
				doneChan <- false
			}

			doneChan <- true
		}()

		select {
		case res := <-doneChan:
			if res {
				log.WithField("task", taskInt).Trace("sending task")
				c <- taskInt
			} else {
				break
			}
		case <-m.stopChan:
			log.Info("got stop signal")
			return
		}
	}
}

// watchNewBlocks subscribes to the best block tracker for new blocks and adds them to the todo list
func (m *Manager) watchNewBlocks() {
	skipBlocks := m.lag
	var lastBlock int64

	for b := range m.tracker.Subscribe() {
		log := log.WithField("block", b)

		if skipBlocks > 0 {
			if b > lastBlock {
				lastBlock = b
				skipBlocks--
			}
			log.Infof("postponing block because lag feature is enabled (%d to go)", skipBlocks)
			continue
		}

		log.Trace("got new block")
		err := m.Todo(b - m.lag)
		if err != nil {
			log.Error(err)
		}

		log.Trace("done adding block to todo")
	}
}
