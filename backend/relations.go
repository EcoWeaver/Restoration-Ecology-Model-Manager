package main

import (
	"fmt"
	"log/slog"
	"slices"
)

type RelationsConfig struct {
	Causal    []string `json:"causal"`
	Noncausal []string `json:"noncausal"`
}

func (self RelationsConfig) verify(l *slog.Logger) error {
	e := map[string]struct{}{}
	for _, s := range []*[]string{&self.Causal, &self.Noncausal} {
		*s = slices.DeleteFunc(*s, func(re string) bool {
			_, ok := e[re]
			e[re] = struct{}{}
			if ok {
				l.Warn(fmt.Sprintf("ignoring duplicate relation kind '%s'", re))
			}
			return ok
		})
	}

	return nil
}
