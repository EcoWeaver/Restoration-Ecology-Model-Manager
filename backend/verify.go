package main

import (
	"backend/types"
	"backend/util"
	"backend/verify"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"net/url"
	"slices"
	"strings"
	"time"

	"golang.org/x/text/language"
)

func (self PutModelRequest) verify(s *State, l *slog.Logger, ctx context.Context) error {
	if err := verifyModel(self.Model, s, l, ctx); err != nil {
		return err
	}

	if self.Parent != nil {
		if err := verify.UUID(*self.Parent); err != nil {
			return err
		}
	}

	return nil
}

func (self PostModelRunRequest) Verify() error {
	type Relation struct {
		From int
		To   int
	}

	if err := verify.UUID(self.UUID); err != nil {
		return err
	}

	e := map[Relation]struct{}{}

	for _, w := range self.RelationWeights {
		pair := Relation{w.From, w.To}

		if _, ok := e[pair]; ok {
			return types.Error{E: fmt.Errorf("encountered duplicate weight for relation from %d to %d", w.From, w.To)}
		}

		e[pair] = struct{}{}
	}

	return nil
}

func modelsEqual(l types.Model, r types.Model) bool {
	return slices.Equal(l.Graph.Nodes, r.Graph.Nodes) && util.SlicesEqualUnordered(l.Graph.Relations, r.Graph.Relations) && util.SlicesEqualUnordered(l.Metadata, r.Metadata)
}

func verifyModel(self types.Model, s *State, l *slog.Logger, ctx context.Context) error {
	if err := verifyGraph(self.Graph, s, l); err != nil {
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
		delete(e, me.Kind)
	}
	if len(e) != 0 {
		var k []string
		for me := range e {
			k = append(k, me)
		}
		return types.Error{E: fmt.Errorf("missing mandatory metadata elements '%s'", strings.Join(k, "', '"))}
	}

	for _, err := range util.SlicesParMap(self.Metadata, func(me types.Metadata) error {
		return verifyMetadata(me, s, m, l, ctx)
	}) {
		if err != nil {
			return err
		}
	}

	return nil
}

func verifyGraph(self types.Graph, s *State, l *slog.Logger) error {
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
		if err := verifyRelation(re, r, len(self.Nodes)); err != nil {
			return err
		}
		pair := Relation{re.From, re.To}
		if _, ok := e[pair]; ok {
			return types.Error{E: fmt.Errorf("encountered duplicate relation from %d to %d", re.From, re.To)}
		}
		e[pair] = struct{}{}
		if slices.Contains(r.Noncausal, re.Kind) {
			pair = Relation{re.To, re.From}
			if _, ok := e[pair]; ok {
				return types.Error{E: fmt.Errorf("encountered duplicate relation from %d to %d", re.From, re.To)}
			}
			e[pair] = struct{}{}
		}
	}

	return nil
}

const Year = "2006"
const YearMonth = "2006-01"
const YearMonthDay = "2006-01-02"

func verifyMetadata(self types.Metadata, s *State, m *MetadataConfig, l *slog.Logger, ctx context.Context) error {
	if self.Kind == "" {
		if self.Value == "" {
			return types.Error{E: errors.New("encountered unset metadata kind with unset value")}
		} else {
			return types.Error{E: fmt.Errorf("encountered unset metadata kind with value '%s'", self.Value)}
		}
	} else if self.Value == "" {
		return types.Error{E: fmt.Errorf("encountered unset value of metadata kind '%s'", self.Kind)}
	}

	for _, ml := range [][]MetadataElement{m.Mandatory, m.Recommended, m.Optional} {
		for _, me := range ml {
			if me.Name != self.Kind {
				continue
			}

			var err error

			switch me.Type {
			case "":
				return nil
			case "choice":
				if slices.Contains(me.Choices, self.Value) {
					return nil
				}
			case "date":
				if _, err := time.Parse(YearMonthDay, self.Value); err == nil {
					return nil
				}
				if _, err := time.Parse(YearMonth, self.Value); err == nil {
					return nil
				}
				if _, err := time.Parse(Year, self.Value); err == nil {
					return nil
				}
			case "doi":
				if err = s.DOICache.Verify(self.Value, ctx); err == nil {
					return nil
				}
			case "ietf-bcp-47":
				if _, err := language.Parse(self.Value); err == nil {
					return nil
				}
			case "ols4":
				ca := verify.OLS4Cache(*me.OLS4)
				c, _ := s.OLS4Caches.LoadOrStore(*me.OLS4, &ca)
				if err = c.(*verify.Cache).Verify(self.Value, ctx); err == nil {
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
			case "url":
				if _, err := url.ParseRequestURI(self.Value); err == nil {
					return nil
				}
			case "year":
				if _, err := time.Parse(Year, self.Value); err == nil {
					return nil
				}
			default:
				l.Warn(fmt.Sprintf("encountered unknown metadata type '%s', assuming free text", me.Type))
				return nil
			}

			if err == nil {
				return types.Error{E: fmt.Errorf("encountered invalid value '%s' of metadata kind '%s', expected type '%s'", self.Value, self.Kind, me.Type)}
			} else if _, ok := err.(types.Error); ok {
				return types.Error{E: fmt.Errorf("encountered invalid value '%s' of metadata kind '%s': %w", self.Value, self.Kind, err)}
			} else {
				return err
			}
		}
	}

	return types.Error{E: fmt.Errorf("encountered invalid metadata kind '%s'", self.Kind)}
}

func verifyRelation(self types.Relation, r *RelationsConfig, n int) error {
	if self.From < 0 || self.From >= n || self.To < 0 || self.To >= n || self.From == self.To {
		return types.Error{E: fmt.Errorf("encountered invalid relation from %d to %d", self.From, self.To)}
	}

	if !slices.Contains(r.Causal, self.Kind) && !slices.Contains(r.Noncausal, self.Kind) {
		return types.Error{E: fmt.Errorf("encountered invalid relation kind '%s'", self.Kind)}
	}

	return verifyPolar(self.Position)
}

func verifyPolar(self types.Polar) error {
	if self.R < 0.0 {
		return types.Error{E: fmt.Errorf("encountered invalid radius '%f'", self.R)}
	}

	if self.Phi < 0.0 || self.Phi >= 2.0*math.Pi {
		return types.Error{E: fmt.Errorf("encountered invalid angle '%f rad'", self.Phi)}
	}

	return nil
}
