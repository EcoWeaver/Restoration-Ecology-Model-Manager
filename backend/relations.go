package main

import (
	"encoding/json"
	"os"
	"slices"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
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

func (self *LiveRelationsConfig) Get(l echo.Logger) (RelationsConfig, error) {
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

	var r RelationsConfig
	err = json.Unmarshal(data, &r)
	if err != nil {
		return RelationsConfig{}, err
	}

	e := map[string]struct{}{}

	for _, s := range []*[]string{&r.Causal, &r.Noncausal} {
		*s = slices.DeleteFunc(*s, func(re string) bool {
			_, ok := e[re]
			e[re] = struct{}{}
			if ok {
				l.Warnf("ignoring duplicate relation kind '%s'", re)
			}
			return ok
		})
	}

	self.config = r
	self.lastMtime = mtime

	return self.config, nil
}
