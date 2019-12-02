package bestblock

import (
	"sync"
	"time"

	"github.com/alethio/web3-go/ethrpc"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("module", "eth")

type Config struct {
	NodeURL      string
	NodeURLWS    string
	PollInterval time.Duration
}

type Tracker struct {
	config          Config
	bestBlockNumber int64
	mu              sync.Mutex
	errChan         chan error
	stopChan        chan bool
	conn            *ethrpc.ETH

	started    bool
	stopped    bool
	subscribed bool

	subscribers map[chan int64]bool
	subsMutex   sync.RWMutex
}

// NewTracker instantiates a new best block tracker with the config provided by the user
// You should call `Run()` to actually start the process
func NewTracker(config Config) (*Tracker, error) {
	return &Tracker{
		config:      config,
		errChan:     make(chan error),
		stopChan:    make(chan bool),
		subscribers: make(map[chan int64]bool),
	}, nil
}

// Run starts the tracker either on websockets (with subscription) or on http (with polling), depending on the
// type of url provided via config
// It takes care of restarting the websocket connection automatically if it crashes
func (b *Tracker) Run() {
	log.Info("starting best block tracker")

	for {
		var conn *ethrpc.ETH
		var err error

		if b.config.NodeURLWS != "" {
			log.Info("setting up websocket connection")
			conn, err = ethrpc.NewWithDefaults(b.config.NodeURLWS)
		} else {
			log.Info("setting up best block polling")
			conn, err = ethrpc.NewWithDefaults(b.config.NodeURL)
		}
		if err != nil {
			log.Error(err)
			time.Sleep(5 * time.Second)
			continue
		}

		b.conn = conn

		if b.started {
			// this happens if, for example, the connection with the node breaks mid-flight
			// we do this to avoid gaps (blocks lost while the connection was down)
			b.getBestHTTP()
		} else {
			// if we're starting the tracker for the first time, we just have to get the current best block and publish it
			b.getBlockNumber()
			b.publish(b.BestBlock())
		}

		if b.config.NodeURLWS != "" {
			b.runWS()
		} else {
			b.runHTTP()
		}

		if b.stopped {
			return
		}
	}
}

// BestBlock returns the current best block known to the tracker
func (b *Tracker) BestBlock() int64 {
	b.mu.Lock()
	block := b.bestBlockNumber
	b.mu.Unlock()

	return block
}

// publish sends the given block to all the clients that are currently subscribed
func (b *Tracker) publish(block int64) {
	b.subsMutex.Lock()
	defer b.subsMutex.Unlock()

	log.WithField("block", block).Tracef("will publish to %d clients", len(b.subscribers))

	for c := range b.subscribers {
		cLocal := c
		// do this to avoid blocking the tracker if the consumer is busy
		go func() {
			cLocal <- block
		}()
	}
}

// Subscribe creates a new channel to send blocks and registers the client on the instance
// returns the channel a client should be consuming from
func (b *Tracker) Subscribe() chan int64 {
	b.subsMutex.Lock()
	defer b.subsMutex.Unlock()

	log.Trace("new client subscribed")

	c := make(chan int64)
	b.subscribers[c] = true
	return c
}

// Err returns a channel of errors that should be consumed to avoid the tracker getting stuck
func (b *Tracker) Err() chan error {
	return b.errChan
}

// Close stops the tracker
func (b *Tracker) Close() {
	if b.subscribed {
		b.stopChan <- true
	}
}
