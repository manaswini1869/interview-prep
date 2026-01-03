package main

import "sync"

type WorkerStats struct {
	WorkerID     string `json:"worker_id"`
	RequestCount int    `json:"request_count"`
}

type Container struct {
	mu       sync.Mutex
	statsMap map[string]int
}

func (c *Container) RecordHit(workerID string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.statsMap[workerID] += 1
}

func (c *Container) GetStats(workerID string) int {
	if c.statsMap == nil || c.statsMap[workerID] == 0 {
		return 0
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.statsMap[workerID]
}
