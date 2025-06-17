package main

import (
	"backend/verify"
	"fmt"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/text/language"
)

type ModelRequest struct {
	Identifier string     `json:"identifier"`
	Graph      Graph      `json:"graph"`
	Version    string     `json:"version"`
	Metadata   []Metadata `json:"metadata"`
}

type Graph struct {
	Nodes     []Node     `json:"nodes"`
	Relations []Relation `json:"relations"`
}

type Metadata struct {
	Kind  string `json:"kind"`
	Value string `json:"value"`
}

type Node struct {
	Body     string `json:"body"`
	Color    string `json:"color"`
	Position Point  `json:"position"`
}

type Relation struct {
	Kind string `json:"kind"`
	From int    `json:"from"`
	To   int    `json:"to"`
}

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func (self *ModelRequest) verify(s *State, l echo.Logger) error {
	if err := s.DoiCache.Verify(self.Identifier); err != nil {
		return err
	}

	if err := verify.Version(self.Version); err != nil {
		return err
	}

	if err := self.Graph.verify(s, l); err != nil {
		return err
	}

	m, err := s.Metadata.Get(l)
	if err != nil {
		return err
	}

	for _, me := range self.Metadata {
		if err = me.verify(s, &m, l); err != nil {
			return err
		}
	}

	return nil
}

func (self Graph) verify(s *State, l echo.Logger) error {
	type Relation struct {
		From int
		To   int
	}

	relations := map[Relation]struct{}{}
	r, err := s.Relations.Get(l)
	if err != nil {
		return err
	}

	for _, relation := range self.Relations {
		if relation.From < 0 || relation.From >= len(self.Nodes) {
			return verify.Error{E: fmt.Errorf("encountered invalid relation from %d", relation.From)}
		}
		if relation.To < 0 || relation.To >= len(self.Nodes) {
			return verify.Error{E: fmt.Errorf("encountered invalid relation to %d", relation.To)}
		}
		pair := Relation{relation.From, relation.To}
		if _, ok := relations[pair]; ok {
			return verify.Error{E: fmt.Errorf("encountered duplicate relation from %d to %d", relation.From, relation.To)}
		}
		relations[pair] = struct{}{}
		if err := relation.verify(&r); err != nil {
			return err
		}
	}

	return nil
}

const Year = "2006"
const YearMonth = "2006-01"
const YearMonthDay = "2006-01-02"

func (self *Metadata) verify(s *State, m *MetadataConfig, l echo.Logger) error {
	self.Value = strings.TrimSpace(self.Value)
	if self.Kind == "" {
		if self.Value == "" {
			return verify.Error{E: fmt.Errorf("encountered unset metadata kind with unset value")}
		} else {
			return verify.Error{E: fmt.Errorf("encountered unset metadata kind with value '%s'", self.Value)}
		}
	} else if self.Value == "" {
		return verify.Error{E: fmt.Errorf("encountered unset value of metadata kind '%s'", self.Kind)}
	}

	for _, ml := range [][]MetadataEntry{m.Mandatory, m.Recommended, m.Optional} {
		for _, me := range ml {
			if me.Name != self.Kind {
				continue
			}

			switch me.Type {
			case "":
				return nil
			case "choice":
				for _, ch := range me.Choices {
					if ch.Label == self.Value {
						return nil
					}
				}
			case "ietf-bcp-47":
				if _, err := language.Parse(self.Value); err == nil {
					return nil
				}
			case "date":
				if _, err := time.Parse(YearMonthDay, self.Value); err == nil {
					return nil
				}
				if _, err := time.Parse(YearMonth, self.Value); err == nil {
					return nil
				}
				fallthrough
			case "year":
				if _, err := time.Parse(Year, self.Value); err == nil {
					return nil
				}
			case "url":
				if _, err := url.ParseRequestURI(self.Value); err == nil {
					return nil
				}
			case "orcid":
				if err := s.OrcidCache.Verify(self.Value); err != nil {
					if _, ok := err.(verify.Error); !ok {
						return err
					}
				} else {
					return nil
				}
			case "species":
				if err := s.SpeciesCache.Verify(self.Value); err != nil {
					if _, ok := err.(verify.Error); !ok {
						return err
					}
				} else {
					return nil
				}
			default:
				l.Warnf("encountered unknown metadata type '%s', assuming free text", me.Type)
				return nil
			}

			return verify.Error{E: fmt.Errorf("encountered invalid value '%s' of metadata kind '%s', expected type '%s'", self.Value, self.Kind, me.Type)}
		}
	}

	return verify.Error{E: fmt.Errorf("encountered invalid metadata kind '%s'", self.Kind)}
}

func (self Relation) verify(r *RelationsConfig) error {
	if slices.Contains(r.Causal, self.Kind) || slices.Contains(r.Noncausal, self.Kind) {
		return nil
	}

	return verify.Error{E: fmt.Errorf("encountered invalid relation kind '%s'", self.Kind)}
}
