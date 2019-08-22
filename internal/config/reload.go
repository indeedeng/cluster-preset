package config

import "time"

// ReloadConfig encapsulates the configurable components of the retry logic.
type ReloadConfig struct {
	FailureRetryInterval time.Duration
	ReloadInterval       time.Duration
}
