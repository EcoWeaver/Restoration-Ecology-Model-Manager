package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
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
	Name    string    `json:"name"`
	Type    *string   `json:"type"`
	Choices *[]Choice `json:"choices"`
}

type Choice struct {
	Label       string  `json:"label"`
	OntologyUri *string `json:"ontology-uri"`
}

func (self *LiveMetadataConfig) Get() (MetadataConfig, error) {
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

	for _, me := range concat(m.Mandatory, m.Recommended, m.Optional) {
		if _, ok := e[me.Name]; ok {
			return MetadataConfig{}, fmt.Errorf("encountered duplicate metadata entry %s", me.Name)
		}
		e[me.Name] = struct{}{}
		err = me.verify()
		if err != nil {
			return MetadataConfig{}, err
		}
	}

	self.config = m
	self.lastMtime = mtime

	return m, nil
}

func (self MetadataEntry) verify() error {
	needsChoices := self.Type != nil && *self.Type == "choice"
	hasChoices := self.Choices != nil && len(*self.Choices) != 0

	if !needsChoices && hasChoices {
		return fmt.Errorf("metadata kind '%s' unexpectedly provides choices", self.Name)
	}

	if needsChoices && !hasChoices {
		return fmt.Errorf("metadata kind '%s' is missing choices", self.Name)
	}

	return nil
}
