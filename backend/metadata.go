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
	Label       string  `json:"label"`
	OntologyUri *string `json:"ontology-uri,omitempty"`
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
	err = json.Unmarshal(data, &m)
	if err != nil {
		return MetadataConfig{}, err
	}

	e := map[string]struct{}{}

	for _, me := range slices.Concat(m.Mandatory, m.Recommended, m.Optional) {
		if _, ok := e[me.Name]; ok {
			l.Warnf("ignoring duplicate metadata kind '%s' of type '%s'", me.Name, me.Type)
		}
		e[me.Name] = struct{}{}
		err = me.verify(l)
		if err != nil {
			return MetadataConfig{}, err
		}
	}

	self.config = m
	self.lastMtime = mtime

	return m, nil
}

func (self *MetadataEntry) verify(l echo.Logger) error {
	// type "choice" needs choices
	// type "" can have choices
	// all other types can't have choices

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
