package config

import (
	"fmt"
	"io/ioutil"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"k8s.io/api/settings/v1alpha1"
)

var (
	errFileRead = fmt.Errorf("failed to read config file")
	errUnmarshal = fmt.Errorf("failed to unmarshal config data")
)

// NewReloadingConfig sets up a new holder that periodically reloads the configuration from disk.
// The reloaded configuration then replaces the current version in memory.
func NewReloadingConfig(path string, reloadConfig *ReloadConfig) (*Holder, error) {
	holder := &Holder{
		path: path,
		mu: &sync.RWMutex{},
		current: nil,
	}

	if err := holder.Reload(); err != nil {
		return nil, err
	}

	go func() {
		for ; true ; {
			sleepDuration := reloadConfig.ReloadInterval
			if err := holder.Reload(); err != nil {
				logrus.Errorf("failed to reload config: %s", err.Error())
				sleepDuration = reloadConfig.FailureRetryInterval
			}
			time.Sleep(sleepDuration)
		}
	}()

	return holder, nil
}

// Holder encapsulates te current configuration.
type Holder struct {
	path string

	mu *sync.RWMutex
	current *v1alpha1.PodPresetSpec
}

// Reload forces a reload of the configuration the holder was configured with.
func (c *Holder) Reload() error {
	data, err := ioutil.ReadFile(c.path)
	if err != nil {
		return errFileRead
	}

	newVersion := &v1alpha1.PodPresetSpec{}
	if err := yaml.Unmarshal(data, newVersion); err != nil {
		return errUnmarshal
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.current = newVersion
	return nil
}

// Get returns the current active version of the configuration.
func (c *Holder) Get() *v1alpha1.PodPresetSpec {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.current
}
