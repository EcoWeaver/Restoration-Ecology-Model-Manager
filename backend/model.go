package main

import (
	"backend/verify"
	"context"
	"errors"
	"fmt"
	"math"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/text/language"
)

type Model struct {
	UUID     string        `json:"uuid"`
	Uploaded time.Time     `json:"uploaded"`
	Graph    Graph         `json:"graph"`
	Parent   *SearchResult `json:"parent"`
	Metadata []Metadata    `json:"metadata"`
}

type UnsavedModel struct {
	Graph    Graph      `json:"graph"`
	Parent   *string    `json:"parent"`
	Metadata []Metadata `json:"metadata"`
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
	Body     string    `json:"body"`
	Color    string    `json:"color"`
	Position Cartesian `json:"position"`
}

type Relation struct {
	Kind     string `json:"kind"`
	From     int    `json:"from"`
	To       int    `json:"to"`
	Position Polar  `json:"position"`
}

type Cartesian struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

type Polar struct {
	R   float32 `json:"r"`
	Phi float32 `json:"phi"`
}

func (self UnsavedModel) verify(s *State, l echo.Logger, ctx context.Context) error {
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
		var k []string
		for me := range e {
			k = append(k, me)
		}
		return verify.Error{E: fmt.Errorf("missing mandatory metadata elements '%s'", strings.Join(k, "', '"))}
	}

	if self.Parent != nil {
		if err := verify.UUID(*self.Parent); err != nil {
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

			var err error

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
			case "doi":
				if err = s.DOICache.Verify(self.Value, ctx); err == nil {
					return nil
				}
			case "orcid":
				if err = s.ORCIDCache.Verify(self.Value, ctx); err == nil {
					return nil
				}
			case "species":
				if err = s.SpeciesCache.Verify(self.Value, ctx); err == nil {
					return nil
				}
			default:
				l.Warnf("encountered unknown metadata type '%s', assuming free text", me.Type)
				return nil
			}

			if err == nil {
				return verify.Error{E: fmt.Errorf("encountered invalid value '%s' of metadata kind '%s', expected type '%s'", self.Value, self.Kind, me.Type)}
			} else if _, ok := err.(verify.Error); ok {
				return verify.Error{E: fmt.Errorf("encountered invalid value '%s' of metadata kind '%s': %w", self.Value, self.Kind, err)}
			} else {
				return err
			}
		}
	}

	return verify.Error{E: fmt.Errorf("encountered invalid metadata kind '%s'", self.Kind)}
}

func (self Relation) verify(r RelationsConfig) error {
	if !slices.Contains(r.Causal, self.Kind) && !slices.Contains(r.Noncausal, self.Kind) {
		return verify.Error{E: fmt.Errorf("encountered invalid relation kind '%s'", self.Kind)}
	}

	return self.Position.verify()
}

func (self Polar) verify() error {
	if self.R < 0.0 {
		return verify.Error{E: fmt.Errorf("encountered invalid radius '%f'", self.R)}
	}

	if self.Phi < 0.0 || self.Phi >= 2.0*math.Pi {
		return verify.Error{E: fmt.Errorf("encountered invalid angle '%f rad'", self.Phi)}
	}

	return nil
}
