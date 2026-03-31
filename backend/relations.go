package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"slices"
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

func (self *LiveRelationsConfig) Get(l *slog.Logger) (RelationsConfig, error) {
	fi, err := os.Stat("relations.json")
	if err != nil {
		return RelationsConfig{}, err
	}

	self.lock.Lock()
	defer self.lock.Unlock()

	mtime := fi.ModTime()
	if time.Time.Equal(mtime, self.lastMtime) {
		return self.config, nil
	}

	data, err := os.ReadFile("relations.json")
	if err != nil {
		return RelationsConfig{}, err
	}

	var r RelationsConfig
	if err = json.Unmarshal(data, &r); err != nil {
		return RelationsConfig{}, err
	}

	e := map[string]struct{}{}
	for _, s := range []*[]string{&r.Causal, &r.Noncausal} {
		*s = slices.DeleteFunc(*s, func(re string) bool {
			_, ok := e[re]
			e[re] = struct{}{}
			if ok {
				l.Warn(fmt.Sprintf("ignoring duplicate relation kind '%s'", re))
			}
			return ok
		})
	}

	self.config = r
	self.lastMtime = mtime

	return self.config, nil
}
