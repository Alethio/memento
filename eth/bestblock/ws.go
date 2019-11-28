package bestblock

import (
	"strconv"
	"time"

	"github.com/alethio/web3-go/types"
)

// runWS uses a websocket subscription to get new block headers as soon as they appear
func (b *Tracker) runWS() {
	heads, err := b.conn.NewHeadsSubscription()
	if err != nil {
		log.Error(err)
		time.Sleep(5 * time.Second)
		return
	}

	b.consumeWSSubscription(heads)

	if b.stopped {
		return
	}

	log.Warn("WS connection closed")
}

// consumeWSSubscription consumes a block headers stream coming from the websocket subscription until the channel is closed
// and saves the new blocks in the dedicated variable on the Tracker struct
func (b *Tracker) consumeWSSubscription(blocks chan *types.BlockHeader) {
	b.subscribed = true
	defer func() { b.subscribed = false }()

	for {
		select {
		case block := <-blocks:
			if block == nil {
				return
			}

			newBlockNumber, _ := strconv.ParseInt(block.Number, 0, 64)

			log.WithField("block", newBlockNumber).Trace("got best block")

			best := b.BestBlock()
			if best >= newBlockNumber {
				// if the new block coming from the node is lower than the current known best
				// we're most likely dealing with a reorg
				for i := newBlockNumber; i <= best; i++ {
					b.publish(i)
				}
			} else {
				for i := best + 1; i <= newBlockNumber; i++ {
					b.publish(i)
				}
			}

			b.mu.Lock()
			b.bestBlockNumber = newBlockNumber
			b.mu.Unlock()
		case <-b.stopChan:
			b.stopped = true
			b.conn.Stop()
			return
		}
	}
}
