package main

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

type LiveMetadataConfig struct {
	config    MetadataConfig
	lock      sync.Mutex
	lastMtime time.Time
}

type MetadataConfig struct {
	Mandatory   []MetadataEntry `json:"mandatory"`
	Recommended []MetadataEntry `json:"recommended"`
	Optional    []MetadataEntry `json:"optional"`
}

type MetadataEntry struct {
	Name    string   `json:"name"`
	Type    string   `json:"type,omitempty"`
	Choices []Choice `json:"choices,omitempty"`
}

type Choice struct {
	Label    string  `json:"label"`
	Ontology *string `json:"ontology,omitempty"`
}

func (self *LiveMetadataConfig) Get(l echo.Logger) (MetadataConfig, error) {
	fi, err := os.Stat("metadata.json")
	if err != nil {
		return MetadataConfig{}, err
	}

	self.lock.Lock()
	defer self.lock.Unlock()

	mtime := fi.ModTime()
	if mtime == self.lastMtime {
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
	for _, s := range []*[]MetadataEntry{&m.Mandatory, &m.Recommended, &m.Optional} {
		*s = slices.DeleteFunc(*s, func(me MetadataEntry) bool {
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
				l.Warnf("ignoring duplicate relation kind '%s'", me)
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

func (self *MetadataEntry) verify(l echo.Logger) error {
	if len(self.Choices) != 0 {
		if self.Type == "" || self.Type == "choice" {
			e := map[string]struct{}{}
			self.Choices = slices.DeleteFunc(self.Choices, func(ce Choice) bool {
				_, ok := e[ce.Label]
				e[ce.Label] = struct{}{}
				if ok {
					l.Warnf("ignoring duplicate choice '%s' for metadata kind '%s'", ce.Label, self.Name)
				}
				return ok
			})
		} else {
			l.Warnf("ignoring unnecessary choices for metadata kind '%s'", self.Name)
			self.Choices = []Choice{}
		}
	} else if self.Type == "choice" {
		return fmt.Errorf("missing choices for metadata kind '%s'", self.Name)
	}

	return nil
}

func (self *MetadataConfig) IsMetadataKind(k string) bool {
	for _, s := range []*[]MetadataEntry{&self.Mandatory, &self.Recommended, &self.Optional} {
		for _, se := range *s {
			if k == se.Name {
				return true
			}
		}
	}

	return false
}
