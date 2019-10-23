package metrics

func (m *Metrics) GetProcessingTime() string {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.processingTime.String()
}

func (m *Metrics) GetScrapingTime() string {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.scrapingTime.String()
}

func (m *Metrics) GetIndexingTime() string {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.indexingTime.String()
}

func (m *Metrics) GetLatestBLock() int64 {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.latestBlock
}

func (m *Metrics) GetTodoLength() int64 {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.todoLength
}

func (m *Metrics) GetReorgedBlocks() int64 {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.reorgedBlocks
}

func (m *Metrics) GetInvalidBlocks() int64 {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.invalidBlocks
}
