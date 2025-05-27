package main

import (
	"fmt"
	"net/url"
	"regexp"
	"slices"
	"time"

	"github.com/docker/go-units"
	"github.com/labstack/echo/v4"
	"golang.org/x/text/language"
)

type ModelRequest struct {
	Identifier string     `json:"identifier"`
	Model      Model      `json:"model"`
	Metadata   []Metadata `json:"metadata"`
}

type Model struct {
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

type GraphVerifyError struct{ err error }

func (err GraphVerifyError) Error() string {
	return err.err.Error()
}

func (self ModelRequest) verify(s *State, l echo.Logger) error {
	err := self.Model.verify(s, l)
	if err != nil {
		return err
	}

	m, err := s.Metadata.Get(l)
	if err != nil {
		return err
	}

	for _, me := range self.Metadata {
		err = me.verify(&m, l)
		if err != nil {
			return err
		}
	}

	return nil
}

func (self Model) verify(s *State, l echo.Logger) error {
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
		if relation.From >= len(self.Nodes) {
			return GraphVerifyError{fmt.Errorf("encountered invalid relation from %d", relation.From)}
		}
		if relation.To >= len(self.Nodes) {
			return GraphVerifyError{fmt.Errorf("encountered invalid relation to %d", relation.To)}
		}
		pair := Relation{relation.From, relation.To}
		if _, ok := relations[pair]; ok {
			return GraphVerifyError{fmt.Errorf("encountered duplicate relation from %d to %d", relation.From, relation.To)}
		}
		relations[pair] = struct{}{}
		err := relation.verify(&r)
		if err != nil {
			return err
		}
	}

	return nil
}

const YearOnly = "2006"

var orcidRegexp = regexp.MustCompile(`\d{4}-\d{4}-\d{4}-\d{4}`)

func (self Metadata) verify(m *MetadataConfig, l echo.Logger) error {
	for _, me := range slices.Concat(m.Mandatory, m.Recommended, m.Optional) {
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
		case "orcid":
			if orcidRegexp.MatchString(self.Value) {
				return nil
			}
		case "url":
			if _, err := url.Parse(self.Value); err == nil {
				return nil
			}
		case "ietf-bcp-47":
			if _, err := language.Parse(self.Value); err == nil {
				return nil
			}
		case "date":
			if _, err := time.Parse(time.DateOnly, self.Value); err == nil {
				return nil
			}
			fallthrough
		case "year":
			if _, err := time.Parse(YearOnly, self.Value); err == nil {
				return nil
			}
		case "bytes-size":
			if _, err := units.RAMInBytes(self.Value); err == nil {
				return nil
			}
		default:
			l.Warnf("encountered unknown metadata kind '%s', assuming 'string'", self.Kind)
			return nil
		}

		return GraphVerifyError{fmt.Errorf("encountered invalid value '%s' of metadata kind '%s', expected type '%s'", self.Value, self.Kind, me.Type)}
	}

	return GraphVerifyError{fmt.Errorf("encountered invalid metadata kind '%s'", self.Kind)}
}

func (self Relation) verify(r *RelationsConfig) error {
	if slices.Contains(r.Causal, self.Kind) || slices.Contains(r.Noncausal, self.Kind) {
		return nil
	}

	return GraphVerifyError{fmt.Errorf("encountered invalid relation kind '%s'", self.Kind)}
}
