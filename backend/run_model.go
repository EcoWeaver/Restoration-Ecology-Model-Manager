package main

import (
	"backend/types"
	"errors"
	"fmt"
	"math"
	"slices"
)

type ModelRunState []*NodeState

type NodeState struct {
	Value    float64
	Change   float64
	In       []Connection
	Out      []*NodeState
	Indegree int
}

type Connection struct {
	From    *NodeState
	Formula types.Polynomial
}

func (self PostModelRunRequest) Build(g types.Graph, r *RelationsConfig) (ModelRunState, error) {
	var s ModelRunState
	for i := range len(g.Nodes) {
		s = append(s, &NodeState{Value: self.NodeValues[i]})
	}

outer:
	for _, w := range self.RelationWeights {
		for _, re := range g.Relations {
			if w.From != re.From || w.To != re.To {
				continue
			}

			if !slices.Contains(r.Causal, re.Kind) {
				return ModelRunState{}, types.Error{E: fmt.Errorf("encountered weight for non-causal relation from %d to %d", w.From, w.To)}
			}

			s[re.From].Out = append(s[re.From].Out, s[re.To])
			s[re.To].In = append(s[re.To].In, Connection{s[re.From], w.Formula})
			s[re.To].Indegree += 1

			continue outer
		}

		return ModelRunState{}, types.Error{E: fmt.Errorf("encountered weight for non-existent relation from %d to %d", w.From, w.To)}
	}

	return s, nil
}

func (self ModelRunState) Run() ([]float64, error) {
	for _, node := range self {
		if node.Indegree == 0 {
			node.recurseKahn()
		}
	}

	var res []float64
	for _, node := range self {
		if node.Indegree != -1 {
			return []float64{}, types.Error{E: errors.New("encountered cycle in graph")}
		}
		res = append(res, node.Value+node.Change)
	}

	return res, nil
}

func (self *NodeState) recurseKahn() {
	self.Indegree -= 1

	for _, in := range self.In {
		if len(in.From.In) == 0 {
			self.Change += eval(in.Formula, in.From.Value)
		} else {
			self.Change += eval(in.Formula, in.From.Change)
		}
	}

	for _, out := range self.Out {
		out.Indegree -= 1
		if out.Indegree == 0 {
			out.recurseKahn()
		}
	}
}

func eval(self types.Polynomial, input float64) float64 {
	var out float64
	for _, coeff := range self {
		out += coeff.Factor * math.Pow(input, coeff.Exponent)
	}
	return out
}
