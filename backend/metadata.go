package main

import (
	"backend/types"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"slices"
	"sync"
	"time"
)

type LiveMetadataConfig struct {
	config    MetadataConfig
	lock      sync.Mutex
	lastMtime time.Time
}

type MetadataConfig struct {
	Mandatory   []MetadataElement `json:"mandatory"`
	Recommended []MetadataElement `json:"recommended"`
	Optional    []MetadataElement `json:"optional"`
}

type MetadataElement struct {
	Name    string      `json:"name"`
	Type    string      `json:"type,omitempty"`
	Choices []string    `json:"choices,omitempty"`
	OLS4    *types.OLS4 `json:"ols4,omitempty"`
}

func (self *LiveMetadataConfig) Get(l *slog.Logger) (MetadataConfig, error) {
	fi, err := os.Stat("metadata.json")
	if err != nil {
		return MetadataConfig{}, err
	}

	self.lock.Lock()
	defer self.lock.Unlock()

	mtime := fi.ModTime()
	if time.Time.Equal(mtime, self.lastMtime) {
		return self.config, nil
	}

	data, err := os.ReadFile("metadata.json")
	if err != nil {
		return MetadataConfig{}, err
	}

	var m MetadataConfig
	if err = json.Unmarshal(data, &m); err != nil {
		return MetadataConfig{}, err
	}

	e := map[string]struct{}{}
	for _, s := range []*[]MetadataElement{&m.Mandatory, &m.Recommended, &m.Optional} {
		*s = slices.DeleteFunc(*s, func(me MetadataElement) bool {
			if err != nil {
				return false
			}
			err = me.verify(l)
			if err != nil {
				return false
			}
			_, ok := e[me.Name]
			e[me.Name] = struct{}{}
			if ok {
				l.Warn(fmt.Sprintf("ignoring duplicate relation kind '%s'", me.Name))
			}
			return ok
		})
		if err != nil {
			return MetadataConfig{}, err
		}
	}

	self.config = m
	self.lastMtime = mtime

	return m, nil
}

func (self *MetadataElement) verify(l *slog.Logger) error {
	if self.Type == "choice" {
		if len(self.Choices) != 0 {
			e := map[string]struct{}{}
			self.Choices = slices.DeleteFunc(self.Choices, func(ce string) bool {
				_, ok := e[ce]
				e[ce] = struct{}{}
				if ok {
					l.Warn(fmt.Sprintf("ignoring duplicate choice '%s' for metadata kind '%s'", ce, self.Name))
				}
				return ok
			})
		} else {
			return fmt.Errorf("missing choices for metadata kind '%s'", self.Name)
		}
	} else if len(self.Choices) != 0 {
		l.Warn(fmt.Sprintf("ignoring unnecessary choices for metadata kind '%s'", self.Name))
		self.Choices = []string{}
	}

	if self.Type == "ols4" {
		if self.OLS4 == nil {
			return fmt.Errorf("missing ontology for metadata kind '%s'", self.Name)
		} else if _, err := url.ParseRequestURI(self.OLS4.IRI); err != nil {
			return err
		}
	} else if self.OLS4 != nil {
		l.Warn(fmt.Sprintf("ignoring unnecessary ontology for metadata kind '%s'", self.Name))
		self.OLS4 = nil
	}

	return nil
}

func (self MetadataConfig) IsMetadataKind(k string) bool {
	return slices.ContainsFunc([][]MetadataElement{self.Mandatory, self.Recommended, self.Optional}, func(ml []MetadataElement) bool {
		return slices.ContainsFunc(ml, func(me MetadataElement) bool {
			return k == me.Name
		})
	})
}
