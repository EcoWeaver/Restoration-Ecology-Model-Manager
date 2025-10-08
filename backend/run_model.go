package main

import (
	"backend/verify"
	"errors"
	"fmt"
	"math"
	"slices"
)

type ModelRunState []*NodeState

type NodeState struct {
	Before   float32
	After    float32
	From     []Connection
	To       []*NodeState
	Indegree int
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

func (self ModelRunRequest) Verify(r *RelationsConfig) error {
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
			return verify.Error{E: fmt.Errorf("encountered duplicate relation weight from %d to %d", w.From, w.To)}
		}

		e[pair] = struct{}{}
	}

	return nil
}

func (self ModelRunRequest) Run(g Graph, r *RelationsConfig) (ModelRunResult, error) {
	var s ModelRunState

	for i := range len(g.Nodes) {
		s = append(s, &NodeState{Before: self.NodeValues[i]})
	}

outer:
	for _, w := range self.RelationWeights {
		for _, v := range g.Relations {
			if w.From == v.From && w.To == v.To {
				if !slices.Contains(r.Causal, v.Kind) {
					return ModelRunResult{}, verify.Error{E: fmt.Errorf("encountered relation weight for non-causal relation from %d to %d", w.From, w.To)}
				}

				s[v.To].From = append(s[v.To].From, Connection{s[v.From], w.Formula, v.Kind == "decreases"})
				s[v.To].Indegree += 1
				s[v.From].To = append(s[v.From].To, s[v.To])

				continue outer
			}
		}

		return ModelRunResult{}, verify.Error{E: fmt.Errorf("encountered non-existent relation weight from %d to %d", w.From, w.To)}
	}

	for _, node := range s {
		if node.Indegree == 0 {
			node.recurseKahn()
		}
	}

	var res ModelRunResult
	for _, node := range s {
		if node.Indegree != -1 {
			return ModelRunResult{}, verify.Error{E: errors.New("encountered cycle in graph")}
		}
		res = append(res, node.After)
	}

	return res, nil
}

func (self *NodeState) recurseKahn() {
	self.Indegree -= 1

	var diff float32
	for _, from := range self.From {
		var res float32
		if len(from.From.From) == 0 {
			res = from.Formula.eval(from.From.Before)
		} else {
			res = from.Formula.eval(from.From.After - from.From.Before)
		}

		if from.Negate {
			diff -= res
		} else {
			diff += res
		}
	}

	self.After = self.Before + diff

	for _, to := range self.To {
		to.Indegree -= 1
		if to.Indegree == 0 {
			to.recurseKahn()
		}
	}
}

func (self Polynomial) eval(input float32) float32 {
	var out float32
	for _, coeff := range self {
		out += coeff.Factor * float32(math.Pow(float64(input), float64(coeff.Exponent)))
	}
	return out
}
