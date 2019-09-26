package bestblock

import "time"

// runHTTP polls the node for the best block number every [config.PollInterval]
func (b *Tracker) runHTTP() {
	for {
		select {
		case <-time.Tick(b.config.PollInterval):
			old := b.BestBlock()
			b.getBlockNumber()

			newBlock := b.BestBlock()
			if old < newBlock {
				for i := old + 1; i <= newBlock; i++ {
					b.publish(i)
				}
			} else {
				b.publish(newBlock)
			}

		case <-b.stopChan:
			return
		}
	}
}
