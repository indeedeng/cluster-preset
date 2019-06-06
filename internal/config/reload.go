package config

import "time"

type ReloadConfig struct {
	FailureRetryInterval time.Duration
	ReloadInterval time.Duration
}
