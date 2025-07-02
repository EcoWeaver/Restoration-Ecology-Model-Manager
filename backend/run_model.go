package main

import (
	"backend/verify"
	"fmt"
	"math"
	"slices"
)

type ModelRunState []*NodeState

type NodeState struct {
	Before float32
	After  float32
	From   []Connection
}

type Connection struct {
	From    *NodeState
	Formula Polynomial
	Negate  bool
}

type ModelRunRequest struct {
	UUID            string           `json:"uuid"`
	NodeValues      map[int]float32  `json:"nodeValues"`
	RelationWeights []RelationWeight `json:"relationWeights"`
}

type ModelRunResult []float32

type RelationWeight struct {
	From    int        `json:"from"`
	To      int        `json:"to"`
	Formula Polynomial `json:"formula"`
}

type Polynomial []Coefficient

type Coefficient struct {
	Factor   float32 `json:"factor"`
	Exponent float32 `json:"exponent"`
}

func (self ModelRunRequest) Verify(g Graph, r *RelationsConfig) error {
	type Relation struct {
		From int
		To   int
	}

	if err := verify.UUID(self.UUID); err != nil {
		return err
	}

	e := map[Relation]struct{}{}
outer:
	for _, rw := range self.RelationWeights {
		pair := Relation{rw.From, rw.To}
		if _, ok := e[pair]; ok {
			return verify.Error{E: fmt.Errorf("encountered duplicate relation weight from %d to %d", rw.From, rw.To)}
		}
		e[pair] = struct{}{}

		for _, re := range g.Relations {
			if re.From == rw.From && re.To == rw.To {
				if slices.Contains(r.Causal, re.Kind) {
					continue outer
				} else {
					return verify.Error{E: fmt.Errorf("encountered relation weight for non-causal relation from %d to %d", rw.From, rw.To)}
				}
			}
		}

		return verify.Error{E: fmt.Errorf("encountered non-existent relation weight from %d to %d", rw.From, rw.To)}
	}

	return nil
}

func (self ModelRunRequest) Run(g Graph, r *RelationsConfig) ModelRunResult {
	var s ModelRunState

	for i := range len(g.Nodes) {
		s = append(s, &NodeState{Before: self.NodeValues[i]})
	}

	for _, v := range g.Relations {
		for _, w := range self.RelationWeights {
			if w.From == v.From && w.To == v.To {
				s[v.To].From = append(s[v.To].From, Connection{s[v.From], w.Formula, v.Kind == "decreases"})
				break
			}
		}
	}

	slices.SortFunc(s, func(a, b *NodeState) int {
		for _, v := range b.From {
			if v.From == a {
				return -1
			}
		}

		return 0
	})

	var res ModelRunResult

	for _, node := range s {
		node.run()
		res = append(res, node.After)
	}

	return res
}

func (self *NodeState) run() {
	self.After = self.Before
	for _, from := range self.From {
		var diff float32
		if len(from.From.From) == 0 {
			diff = from.Formula.run(from.From.Before)
		} else {
			diff = from.Formula.run(from.From.After - from.From.Before)
		}
		if from.Negate {
			diff *= -1
		}
		self.After += diff
	}
}

func (self Polynomial) run(input float32) float32 {
	var out float32
	for _, coeff := range self {
		out += coeff.Factor * float32(math.Pow(float64(input), float64(coeff.Exponent)))
	}
	return out
}
