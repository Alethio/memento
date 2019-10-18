package core

import (
	"database/sql"
	"sync"
	"time"

	"github.com/pressly/goose"

	"git.aleth.io/alethio/memento/api"

	"github.com/alethio/web3-go/validator"

	"git.aleth.io/alethio/memento/scraper"
	"git.aleth.io/alethio/memento/taskmanager"

	"git.aleth.io/alethio/memento/eth/bestblock"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("module", "core")

type Core struct {
	config Config

	bbtracker   *bestblock.Tracker
	taskmanager *taskmanager.Manager
	scraper     *scraper.Scraper
	db          *sql.DB
	api         *api.API

	stopMu sync.Mutex
}

func New(config Config) *Core {
	bbtracker, err := bestblock.NewTracker(config.BestBlockTracker)
	if err != nil {
		log.Fatal("could not start best block tracker")
		return nil
	}

	go bbtracker.Run()
	err = <-bbtracker.Err()
	if err != nil {
		log.Fatal("could not start best block tracker")
	}

	go func() {
		// todo: can we handle these errors?
		for err := range bbtracker.Err() {
			log.Error(err)
		}
	}()

	var lag int64
	if config.Features.Lag.Enabled {
		lag = config.Features.Lag.Value
	}

	tm, err := taskmanager.New(bbtracker, lag, config.TaskManager)
	if err != nil {
		log.Fatal("could not start task manager")
	}

	s, err := scraper.New(config.Scraper)
	if err != nil {
		log.Fatal("could not start scraper")
	}

	log.Info("connecting to postgres")
	db, err := sql.Open("postgres", config.PostgresConnectionString)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	if config.Features.Automigrate {
		log.Info("attempting automatic execution of migrations")
		err = goose.Up(db, "/")
		if err != nil && err != goose.ErrNoNextVersion {
			log.Fatal(err)
		}
		log.Info("database version is up to date")
	}

	log.Info("connected to postgres successfuly")

	a := api.New(db, config.API)

	return &Core{
		config:      config,
		bbtracker:   bbtracker,
		taskmanager: tm,
		scraper:     s,
		db:          db,
		api:         a,
	}
}

func (c *Core) Run() {
	go c.api.Start()

	blockChan := make(chan int64)

	max, err := c.getHighestBlock()
	if err != nil {
		log.Fatal("could not get highest block from db:", err)
	}

	log.WithField("block", max).Info("got highest block from db")

	best := c.bbtracker.BestBlock()

	log.WithField("block", best).Info("got highest block from network")

	if c.config.Features.Backfill {
		var lag int64
		if c.config.Features.Lag.Enabled {
			lag = c.config.Features.Lag.Value
		}

		backfillTarget := best - lag

		if max+1 < backfillTarget {
			log.Infof("adding tasks for %d blocks to be backfilled", backfillTarget-max+1)
			for i := max; i <= backfillTarget; i++ {
				err := c.taskmanager.Todo(i)
				if err != nil {
					log.Fatal("could not add task:", err)
				}
			}
		}
	} else {
		log.Info("skipping backfilling since feature is disabled")
	}

	go c.taskmanager.FeedToChan(blockChan)

	go func() {
		for b := range blockChan {
			c.stopMu.Lock()
			log := log.WithField("block", b)
			log.Info("processing block")

			start := time.Now()
			blk, err := c.scraper.Exec(b)
			if err != nil {
				c.stopMu.Unlock()
				err1 := c.taskmanager.Todo(b)
				if err1 != nil {
					log.Fatal(err1)
				}
				time.Sleep(2 * time.Second)
				continue
			}

			log.Debug("validating block")
			v := validator.New()
			v.LoadBlock(blk.Block)
			v.LoadUncles(blk.Uncles)
			v.LoadReceipts(blk.Receipts)

			_, err = v.Run()
			if err != nil {
				c.stopMu.Unlock()
				log.Error("error validating block: ", err)
				err1 := c.taskmanager.Todo(b)
				if err1 != nil {
					log.Fatal(err1)
				}
				continue
			}
			log.Debug("block is valid")

			log.Debug("storing block into the database")
			blk.RegisterStorables()
			err = blk.Store(c.db)
			if err != nil {
				c.stopMu.Unlock()
				log.Error("error storing block: ", err)
				err1 := c.taskmanager.Todo(b)
				if err1 != nil {
					log.Fatal(err1)
				}
				continue
			}
			log.WithField("duration", time.Since(start)).Info("done processing block")
			c.stopMu.Unlock()
		}
	}()
}

func (c *Core) Close() error {
	c.stopMu.Lock()
	defer c.stopMu.Unlock()

	c.bbtracker.Close()
	log.Info("closed best block tracker")

	err := c.db.Close()
	if err != nil {
		return err
	}
	log.Info("closed db connection")

	errChan := make(chan error)
	go func() {
		err = c.taskmanager.Close()
		if err != nil {
			errChan <- err
		}
		log.Info("closed task manager")
		errChan <- nil
	}()

	select {
	case err := <-errChan:
		return err
	case <-time.After(5 * time.Second):
		log.Warn("could not close task manager, exiting uncleanly")
	}

	return nil
}
