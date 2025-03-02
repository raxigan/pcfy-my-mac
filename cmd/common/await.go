package common

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var DefaultPollInterval = 100 * time.Millisecond

var DefaultTimeout = 10 * time.Second

type Config struct {
	PollInterval time.Duration
	Timeout      time.Duration
}

func DefaultConfig() Config {
	return Config{
		PollInterval: DefaultPollInterval,
		Timeout:      DefaultTimeout,
	}
}

var ErrTimeout = errors.New("condition not satisfied within timeout period")

type Condition func() bool

func Until(condition Condition, options ...func(*Config)) error {
	cfg := DefaultConfig()
	for _, option := range options {
		option(&cfg)
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	ticker := time.NewTicker(cfg.PollInterval)
	defer ticker.Stop()

	if condition() {
		return nil
	}

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("%w after %v", ErrTimeout, cfg.Timeout)
		case <-ticker.C:
			if condition() {
				return nil
			}
		}
	}
}

func WithPollInterval(interval time.Duration) func(*Config) {
	return func(c *Config) {
		c.PollInterval = interval
	}
}

func WithTimeout(timeout time.Duration) func(*Config) {
	return func(c *Config) {
		c.Timeout = timeout
	}
}
