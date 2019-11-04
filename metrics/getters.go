package metrics

func (p *Provider) GetProcessingTime() string {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.processingTime.String()
}

func (p *Provider) GetRawProcessingTime() int64 {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.processingTime.Raw()
}

func (p *Provider) GetScrapingTime() string {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.scrapingTime.String()
}

func (p *Provider) GetRawScrapingTime() int64 {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.scrapingTime.Raw()
}

func (p *Provider) GetIndexingTime() string {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.indexingTime.String()
}

func (p *Provider) GetRawIndexingTime() int64 {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.indexingTime.Raw()
}

func (p *Provider) GetLatestBLock() int64 {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.latestBlock
}

func (p *Provider) GetTodoLength() int64 {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.todoLength
}

func (p *Provider) GetReorgedBlocks() int64 {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.reorgedBlocks
}

func (p *Provider) GetInvalidBlocks() int64 {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.invalidBlocks
}
