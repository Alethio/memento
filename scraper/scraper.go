package scraper

import (
	"strconv"
	"time"

	"github.com/Alethio/memento/data"

	"github.com/alethio/web3-go/ethrpc"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("module", "scraper")

type Config struct {
	NodeURL      string
	EnableUncles bool
}

type Scraper struct {
	config Config

	conn *ethrpc.ETH
}

func New(config Config) (*Scraper, error) {
	c, err := ethrpc.NewWithDefaults(config.NodeURL)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &Scraper{
		config: config,
		conn:   c,
	}, nil
}

// Exec does the JSONRPC calls necessary for scraping a given block and returns the raw data
// It:
// - scrapes the block using eth_getBlockByNumber
// - for each transaction in the block, scrapes the receipts using eth_getTransactionReceipt
// - for each uncle in the block, scrapes the data using eth_getUncleByBlockHashAndIndex
func (s *Scraper) Exec(block int64) (*data.FullBlock, error) {
	log = log.WithField("block", block)

	b := &data.FullBlock{}

	log.Debug("getting block")
	start := time.Now()
	dataBlock, err := s.conn.GetBlockByNumber("0x" + strconv.FormatInt(block, 16))
	if err != nil {
		log.Error(err)
		return nil, err
	}
	b.Block = dataBlock
	log.WithField("duration", time.Since(start)).Debug("got block")

	log.Debug("getting receipts")
	start = time.Now()
	for _, tx := range dataBlock.Transactions {
		dataReceipt, err := s.conn.GetTransactionReceipt(tx.Hash)
		if err != nil {
			log.Error(err)
			return nil, err
		}

		b.Receipts = append(b.Receipts, dataReceipt)
	}
	log.WithField("duration", time.Since(start)).Debugf("got %d receipts", len(b.Receipts))

	if s.config.EnableUncles {
		log.Debug("getting uncles")
		start = time.Now()
		for idx := range dataBlock.Uncles {
			dataUncle, err := s.conn.GetUncleByBlockHashAndIndex(b.Block.Hash, "0x"+strconv.FormatInt(int64(idx), 16))
			if err != nil {
				log.Error(err)
				return nil, err
			}

			b.Uncles = append(b.Uncles, dataUncle)
		}
		log.WithField("duration", time.Since(start)).Debugf("got %d uncles", len(b.Uncles))
	}

	log.Debug("done scraping block")

	return b, nil
}
