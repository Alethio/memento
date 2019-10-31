package metrics

import (
	"sync"
	"time"
)

type Provider struct {
	mu sync.Mutex

	processingTime AverageDuration
	scrapingTime   AverageDuration
	indexingTime   AverageDuration

	latestBlock   int64
	todoLength    int64
	reorgedBlocks int64
	invalidBlocks int64
}

func New() *Provider {
	return &Provider{}
}

func (p *Provider) Reset() {
	p.processingTime.Reset()
	p.scrapingTime.Reset()
	p.indexingTime.Reset()

	p.latestBlock = 0
	p.todoLength = 0
	p.reorgedBlocks = 0
	p.invalidBlocks = 0
}

func (p *Provider) RecordProcessingTime(duration time.Duration) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.processingTime.Add(duration)
}

func (p *Provider) RecordScrapingTime(duration time.Duration) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.scrapingTime.Add(duration)
}

func (p *Provider) RecordIndexingTime(duration time.Duration) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.indexingTime.Add(duration)
}

func (p *Provider) RecordLatestBlock(block int64) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.latestBlock = block
}

func (p *Provider) RecordTodoLength(len int64) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.todoLength = len
}

func (p *Provider) RecordReorgedBlock() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.reorgedBlocks++
}

func (p *Provider) RecordInvalidBlock() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.invalidBlocks++
}
