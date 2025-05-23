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

func (self ModelRequest) verify(c *Context) error {
	err := self.Model.verify(c)
	if err != nil {
		return err
	}

	m, err := c.Metadata.Get()
	if err != nil {
		return err
	}

	for _, me := range self.Metadata {
		err = me.verify(&m, c.Logger())
		if err != nil {
			return err
		}
	}

	return nil
}

func (self Model) verify(c *Context) error {
	type Relation struct {
		From int
		To   int
	}

	relations := map[Relation]struct{}{}
	r, err := c.Relations.Get()
	if err != nil {
		return err
	}

	for _, relation := range self.Relations {
		if relation.From >= len(self.Nodes) {
			return fmt.Errorf("encountered invalid relation connection from %d", relation.From)
		}
		if relation.To >= len(self.Nodes) {
			return fmt.Errorf("encountered invalid relation connection to %d", relation.To)
		}
		pair := Relation{relation.From, relation.To}
		if _, ok := relations[pair]; ok {
			return fmt.Errorf("encountered duplicate relation from %d to %d", relation.From, relation.To)
		}
		relations[pair] = struct{}{}
		err := relation.verify(&r)
		if err != nil {
			return err
		}
	}

	return nil
}

var orcidRegexp = regexp.MustCompile(`\d{4}-\d{4}-\d{4}-\d{4}`)

func (self Metadata) verify(m *MetadataConfig, l echo.Logger) error {
	for _, me := range concat(m.Mandatory, m.Recommended, m.Optional) {
		if me.Name != self.Kind {
			continue
		}

		if me.Type == nil {
			return nil
		}

		switch *me.Type {
		case "choice":
			for _, ch := range *me.Choices {
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
		case "ietf-rfc-3339":
			if _, err := time.Parse(time.RFC3339, self.Value); err == nil {
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

		return fmt.Errorf("encountered invalid value '%s' of metadata kind '%s'", self.Value, self.Kind)
	}

	return fmt.Errorf("encountered invalid metadata kind '%s'", self.Kind)
}

func (self Relation) verify(r *RelationsConfig) error {
	if slices.Contains(r.Causal, self.Kind) || slices.Contains(r.Noncausal, self.Kind) {
		return nil
	}

	return fmt.Errorf("encountered invalid relation kind '%s'", self.Kind)
}
