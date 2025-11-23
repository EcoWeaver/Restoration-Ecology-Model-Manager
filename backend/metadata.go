package main

import (
	"backend/verify"
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
	Name    string       `json:"name"`
	Type    string       `json:"type,omitempty"`
	Choices []string     `json:"choices,omitempty"`
	OLS4    *verify.OLS4 `json:"ols4,omitempty"`
}

func (self *LiveMetadataConfig) Get(l echo.Logger) (MetadataConfig, error) {
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
	if self.Type == "choice" {
		if len(self.Choices) != 0 {
			e := map[string]struct{}{}
			self.Choices = slices.DeleteFunc(self.Choices, func(ce string) bool {
				_, ok := e[ce]
				e[ce] = struct{}{}
				if ok {
					l.Warnf("ignoring duplicate choice '%s' for metadata kind '%s'", ce, self.Name)
				}
				return ok
			})
		} else {
			return fmt.Errorf("missing choices for metadata kind '%s'", self.Name)
		}
	} else if len(self.Choices) != 0 {
		l.Warnf("ignoring unnecessary choices for metadata kind '%s'", self.Name)
		self.Choices = []string{}
	}

	if self.Type == "ols4" {
		if self.OLS4 == nil {
			return fmt.Errorf("missing ontology for metadata kind '%s'", self.Name)
		}
	} else if self.OLS4 != nil {
		l.Warnf("ignoring unnecessary ontology for metadata kind '%s'", self.Name)
		self.OLS4 = nil
	}

	return nil
}

func (self MetadataConfig) IsMetadataKind(k string) bool {
	return slices.ContainsFunc([][]MetadataEntry{self.Mandatory, self.Recommended, self.Optional}, func(ml []MetadataEntry) bool {
		return slices.ContainsFunc(ml, func(me MetadataEntry) bool {
			return k == me.Name
		})
	})
}
