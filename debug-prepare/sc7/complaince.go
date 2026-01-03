package main

import (
	"context"
	"errors"
	"net/http"
	"time"
)

type ComplianceClient struct {
	Endpoint string
	Client   *http.Client
}

func (c *ComplianceClient) CheckEndpointHealth() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Millisecond)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, "GET", c.Endpoint+"/health", nil)
	resp, err := c.Client.Do(req)
	return err == nil && resp.StatusCode == http.StatusOK
}

func (c *ComplianceClient) HealthHandler(w *http.ResponseWriter, r *http.Request) {
	if !c.CheckEndpointHealth() {
		http.Error(*w, "Compliance service unhealthy", http.StatusServiceUnavailable)
		return
	}
}

func (c *ComplianceClient) CheckCode(ctx context.Context, code string) error {
	// Current implementation: No timeout, no retries.
	// This blocks the entire goroutine if the external service is slow.
	for attempt := 0; attempt < 2; attempt++ {
		if attempt > 0 {
			time.Sleep(100 * time.Millisecond)
		}
		reqCtx, cancel := context.WithTimeout(ctx, 300*time.Millisecond)
		req, _ := http.NewRequestWithContext(reqCtx, "POST", c.Endpoint, nil)
		resp, err := c.Client.Do(req)
		cancel()
		if err != nil {
			continue
		}
		if resp.StatusCode != http.StatusOK {
			continue
		}
		return nil
	}
	return errors.New("compliance check failed after retries")
}
