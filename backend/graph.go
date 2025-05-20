package main

import (
	"backend/metadata"
	"backend/relations"
)

type GraphRequest struct {
	Metadata []Metadata `json:"metadata"`
	Graph    Graph      `json:"graph"`
}

type Graph struct {
	Nodes     []Node     `json:"nodes"`
	Relations []Relation `json:"relations"`
}

type Metadata struct {
	Kind     metadata.Element `json:"kind"`
	Value    string           `json:"value"`
	Ontology string           `json:"ontology"`
}

type Node struct {
	Body     string `json:"body"`
	Color    string `json:"color"`
	Position Point  `json:"position"`
}

type Relation struct {
	From int               `json:"from"`
	To   int               `json:"to"`
	Kind relations.Element `json:"kind"`
}

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}
