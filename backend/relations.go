package main

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

type LiveRelationsConfig struct {
	config    RelationsConfig
	lock      sync.Mutex
	lastMtime time.Time
}

type RelationsConfig struct {
	Causal    []string `json:"causal"`
	Noncausal []string `json:"noncausal"`
}

func (self *LiveRelationsConfig) Get() (RelationsConfig, error) {
	fi, err := os.Stat("relations.json")
	if err != nil {
		return RelationsConfig{}, err
	}

	self.lock.Lock()
	defer self.lock.Unlock()

	mtime := fi.ModTime()
	if mtime == self.lastMtime {
		return self.config, nil
	}

	data, err := os.ReadFile("relations.json")
	if err != nil {
		return RelationsConfig{}, err
	}

	err = json.Unmarshal(data, &self.config)
	if err != nil {
		return RelationsConfig{}, err
	}
	self.lastMtime = mtime

	return self.config, nil
}
