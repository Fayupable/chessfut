package health

import (
	"context"
	"log/slog"
	"sync/atomic"
	"time"
)

var retryDelays = []time.Duration{
	1 * time.Second,
	5 * time.Second,
	15 * time.Second,
	30 * time.Second,
	60 * time.Second,
	5 * time.Minute,
	15 * time.Minute,
}

const healthyCheckInterval = 10 * time.Minute

type Checker struct {
	name    string
	checkFn func(ctx context.Context) error

	available      atomic.Bool
	failedAttempts atomic.Int32
}

func NewChecker(name string, checkFn func(ctx context.Context) error) *Checker {
	return &Checker{name: name, checkFn: checkFn}
}

func (c *Checker) Start(ctx context.Context) {
	go c.loop(ctx)
}

func (c *Checker) loop(ctx context.Context) {
	c.check(ctx)
	for {
		timer := time.NewTimer(c.nextDelay())
		select {
		case <-ctx.Done():
			timer.Stop()
			return
		case <-timer.C:
			c.check(ctx)
		}
	}
}

func (c *Checker) check(ctx context.Context) {
	if err := c.checkFn(ctx); err != nil {
		c.onFailure(err)
		return
	}
	c.onSuccess()
}

func (c *Checker) onSuccess() {
	wasDown := !c.available.Load()
	attempts := c.failedAttempts.Load()

	c.available.Store(true)
	c.failedAttempts.Store(0)

	if wasDown && attempts > 0 {
		slog.Info("health_check_restored", "service", c.name, "after_failed_attempts", attempts)
	}
}

func (c *Checker) onFailure(err error) {
	attempts := c.failedAttempts.Add(1)
	c.available.Store(false)

	slog.Warn("health_check_failed",
		"service", c.name,
		"attempt", attempts,
		"next_retry", c.nextDelay().String(),
		"error", err.Error(),
	)
}

func (c *Checker) nextDelay() time.Duration {
	attempts := int(c.failedAttempts.Load())
	if attempts == 0 {
		return healthyCheckInterval
	}

	idx := attempts - 1
	if idx >= len(retryDelays) {
		idx = len(retryDelays) - 1
	}
	return retryDelays[idx]
}

func (c *Checker) IsAvailable() bool {
	return c.available.Load()
}

func (c *Checker) FailedAttempts() int {
	return int(c.failedAttempts.Load())
}
