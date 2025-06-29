package main

import (
	"backend/verify"
	"context"
	"errors"
	"fmt"
	"net/url"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/text/language"
)

type Model struct {
	UUID     *uuid.UUID `json:"uuid"`
	Graph    Graph      `json:"graph"`
	Metadata []Metadata `json:"metadata"`
	Parent   *uuid.UUID `json:"parent"`
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

func (self Model) verify(s *State, l echo.Logger, ctx context.Context) error {
	if err := self.Graph.verify(s, l); err != nil {
		return err
	}

	m, err := s.Metadata.Get(l)
	if err != nil {
		return err
	}

	e := map[string]struct{}{}
	for _, me := range m.Mandatory {
		e[me.Name] = struct{}{}
	}
	for _, me := range self.Metadata {
		if err = me.verify(s, m, l, ctx); err != nil {
			return err
		}
		delete(e, me.Kind)
	}
	if len(e) != 0 {
		err := "missing mandatory metadata elements"
		for me := range e {
			err += " '" + me + "'"
		}
		return verify.Error{E: errors.New(err)}
	}

	return nil
}

func (self Graph) verify(s *State, l echo.Logger) error {
	type Relation struct {
		From int
		To   int
	}

	e := map[Relation]struct{}{}
	r, err := s.Relations.Get(l)
	if err != nil {
		return err
	}

	for _, ne := range self.Nodes {
		if err := verify.Color(ne.Color); err != nil {
			return err
		}
	}

	for _, re := range self.Relations {
		if re.From < 0 || re.From >= len(self.Nodes) {
			return verify.Error{E: fmt.Errorf("encountered invalid relation from %d", re.From)}
		}
		if re.To < 0 || re.To >= len(self.Nodes) {
			return verify.Error{E: fmt.Errorf("encountered invalid relation to %d", re.To)}
		}
		if re.From == re.To {
			return verify.Error{E: fmt.Errorf("encountered invalid relation from %d to %d", re.From, re.To)}
		}
		pair := Relation{re.From, re.To}
		if _, ok := e[pair]; ok {
			return verify.Error{E: fmt.Errorf("encountered duplicate relation from %d to %d", re.From, re.To)}
		}
		e[pair] = struct{}{}
		if err := re.verify(r); err != nil {
			return err
		}
	}

	return nil
}

const Year = "2006"
const YearMonth = "2006-01"
const YearMonthDay = "2006-01-02"

func (self Metadata) verify(s *State, m MetadataConfig, l echo.Logger, ctx context.Context) error {
	if self.Kind == "" {
		if self.Value == "" {
			return verify.Error{E: errors.New("encountered unset metadata kind with unset value")}
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
				if err := verify.ORCID(self.Value, ctx); err != nil {
					if _, ok := err.(verify.Error); !ok {
						return err
					}
				} else {
					return nil
				}
			case "doi":
				if err := s.DOICache.Verify(self.Value, ctx); err != nil {
					if _, ok := err.(verify.Error); !ok {
						return err
					}
				} else {
					return nil
				}
			case "species":
				if err := s.SpeciesCache.Verify(self.Value, ctx); err != nil {
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

func (self Relation) verify(r RelationsConfig) error {
	if slices.Contains(r.Causal, self.Kind) || slices.Contains(r.Noncausal, self.Kind) {
		return nil
	}

	return verify.Error{E: fmt.Errorf("encountered invalid relation kind '%s'", self.Kind)}
}
