package bestblock

import "time"

// runHTTP polls the node for the best block number every [config.PollInterval]
func (b *Tracker) runHTTP() {
	log.Tracef("tracking best block via HTTP polling, every %s", b.config.PollInterval)
	for {
		select {
		case <-time.Tick(b.config.PollInterval):
			b.getBestHTTP()
		case <-b.stopChan:
			return
		}
	}
}

// getBestHTTP checks the current best block and publishes events for all the blocks between the previous and the new best block
func (b *Tracker) getBestHTTP() {
	oldBest := b.BestBlock()
	b.getBlockNumber()
	newBest := b.BestBlock()

	for i := oldBest + 1; i <= newBest; i++ {
		b.publish(i)
	}
}
