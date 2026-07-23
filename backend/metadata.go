package main

import (
	"backend/types"
	"fmt"
	"log/slog"
	"net/url"
	"slices"
)

type MetadataConfig struct {
	Mandatory   []MetadataElement `json:"mandatory"`
	Recommended []MetadataElement `json:"recommended"`
	Optional    []MetadataElement `json:"optional"`
}

type MetadataElement struct {
	Name    string      `json:"name"`
	Type    string      `json:"type,omitempty"`
	Choices []string    `json:"choices,omitempty"`
	Onto    *types.Onto `json:"onto,omitempty"`
}

func (self MetadataConfig) verify(l *slog.Logger) error {
	var err error

	e := map[string]struct{}{}
	for _, s := range []*[]MetadataElement{&self.Mandatory, &self.Recommended, &self.Optional} {
		*s = slices.DeleteFunc(*s, func(me MetadataElement) bool {
			err = me.verify(l)
			_, ok := e[me.Name]
			e[me.Name] = struct{}{}
			if ok {
				l.Warn(fmt.Sprintf("ignoring duplicate metadata kind '%s'", me.Name))
			}
			return ok
		})
	}

	return err
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
		if self.Onto == nil {
			return fmt.Errorf("missing ontology for metadata kind '%s'", self.Name)
		} else if _, err := url.ParseRequestURI(self.Onto.IRI); err != nil {
			return err
		}
	} else if self.Onto != nil {
		l.Warn(fmt.Sprintf("ignoring unnecessary ontology for metadata kind '%s'", self.Name))
		self.Onto = nil
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
