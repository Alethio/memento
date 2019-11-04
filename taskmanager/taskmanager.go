package taskmanager

import (
	"math"
	"strconv"
	"time"

	"github.com/Alethio/memento/metrics"

	"github.com/sirupsen/logrus"

	"github.com/Alethio/memento/eth/bestblock"
	"github.com/go-redis/redis"
)

var log = logrus.WithField("module", "taskmanager")

type Config struct {
	RedisServer     string
	RedisPassword   string
	TodoList        string
	BackfillEnabled bool
}

type Manager struct {
	config Config
	lag    int64

	metrics *metrics.Provider
	tracker *bestblock.Tracker
	redis   *redis.Client

	paused        bool
	pause, resume chan bool

	lastBlockAdded int64

	closed   bool
	stopChan chan bool
}

// New instantiates a new task manager and also takes care of the redis connection management
// it subscribes to the best block tracker for new blocks which it'll add to the redis queue automatically
func New(tracker *bestblock.Tracker, lag int64, metrics *metrics.Provider, config Config) (*Manager, error) {
	m := &Manager{
		config:   config,
		lag:      lag,
		metrics:  metrics,
		tracker:  tracker,
		closed:   false,
		stopChan: make(chan bool),
		pause:    make(chan bool),
		resume:   make(chan bool),
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
	if !m.paused {
		m.stopChan <- true
	}
	return m.redis.Close()
}

func (m *Manager) Pause() {
	if !m.paused {
		log.Trace("attempting task manager pause")

		// add a member with value "-1" to the queue in order to unblock the BZPopMax and pause the manager completely
		err := m.redis.ZAdd(m.config.TodoList, redis.Z{
			Score:  math.MaxFloat64,
			Member: "-1",
		}).Err()
		if err != nil {
			log.Error(err)
			return
		}

		log.Trace("sending pause signal")
		m.pause <- true
		log.Trace("sent pause signal")
	}
}

func (m *Manager) Resume() {
	if m.paused {
		log.Trace("sending resume signal")
		m.resume <- true
		log.Trace("sent resume signal")

		// need to sync with the FeedToChan function
		// make sure it has enough time to process the signal
		time.Sleep(100 * time.Millisecond)
	}
}

func (m *Manager) IsPaused() bool {
	return m.paused
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

		select {
		case <-m.pause:
			m.paused = true
			log.Trace("task manager is paused")
			<-m.resume
			log.Trace("task manager has resumed")
			m.paused = false
		default:
			break
		}

		// update the metrics with the queue length
		todoLen, err := m.redis.ZCard(m.config.TodoList).Result()
		if err != nil {
			log.Error(err)
		} else {
			m.metrics.RecordTodoLength(todoLen)
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
				// don't send any task on the channel if we got the pause signal
				if taskInt == -1 {
					log.Trace("got pause signal from redis")

					// need to sync with the Pause function
					// make sure it has enough time to send the signal
					time.Sleep(100 * time.Millisecond)
					break
				}

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
	var started bool

	for b := range m.tracker.Subscribe() {
		log := log.WithField("block", b)

		if !started || !m.config.BackfillEnabled || b-m.lag <= m.lastBlockAdded {
			started = true
			m.lastBlockAdded = b - m.lag - 1
		}

		if skipBlocks > 0 {
			if b > lastBlock {
				lastBlock = b
				skipBlocks--
			}
			log.Infof("postponing block because lag feature is enabled (%d to go)", skipBlocks)
			continue
		}

		if m.paused {
			log.Info("skipping block because I'm paused")
			continue
		}

		log.Trace("got new block")

		for i := m.lastBlockAdded + 1; i <= b-m.lag; i++ {
			err := m.Todo(i)
			if err != nil {
				log.Error(err)
			} else {
				m.lastBlockAdded = i
			}
		}

		log.Trace("done adding block to todo")
	}
}
