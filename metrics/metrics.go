package metrics

import (
	"sync"
	"time"
)

type Metrics struct {
	mu sync.Mutex

	processingTime AverageDuration
	scrapingTime   AverageDuration
	indexingTime   AverageDuration

	latestBlock   int64
	todoLength    int64
	reorgedBlocks int64
	invalidBlocks int64
}

func New() *Metrics {
	return &Metrics{}
}

func (m *Metrics) RecordProcessingTime(duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.processingTime.Add(duration)
}

func (m *Metrics) RecordScrapingTime(duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.scrapingTime.Add(duration)
}

func (m *Metrics) RecordIndexingTime(duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.indexingTime.Add(duration)
}

func (m *Metrics) RecordLatestBlock(block int64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.latestBlock = block
}

func (m *Metrics) RecordTodoLength(len int64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.todoLength = len
}

func (m *Metrics) RecordReorgedBlock() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.reorgedBlocks++
}

func (m *Metrics) RecordInvalidBlock() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.invalidBlocks++
}
