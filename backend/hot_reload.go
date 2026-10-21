package main

import (
	"encoding/json"
	"log/slog"
	"os"
	"sync"
	"time"
)

type Verify interface {
	verify(l *slog.Logger) error
}

type HotReload[T Verify] struct {
	path      string
	cached    T
	lock      sync.Mutex
	lastMtime time.Time
}

func (self *HotReload[T]) Get(l *slog.Logger) (*T, error) {
	self.lock.Lock()
	defer self.lock.Unlock()

	if time.Since(self.lastMtime) < time.Second {
		return &self.cached, nil
	}

	fi, err := os.Stat(self.path)
	if err != nil {
		return nil, err
	}

	mtime := fi.ModTime()
	if time.Time.Equal(mtime, self.lastMtime) {
		return &self.cached, nil
	}

	data, err := os.ReadFile(self.path)
	if err != nil {
		return nil, err
	}

	var new T
	if err = json.Unmarshal(data, &new); err != nil {
		return nil, err
	}

	if err := new.verify(l); err != nil {
		return nil, err
	}

	self.cached = new
	self.lastMtime = mtime

	return &self.cached, nil
}
