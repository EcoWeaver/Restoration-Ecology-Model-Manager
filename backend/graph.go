package main

import (
	"backend/metadata"
	"backend/relations"
)

type Graph struct {
	Metadata  []Metadata `json:"metadata"`
	Nodes     []Node     `json:"nodes"`
	Relations []Relation `json:"relations"`
}

type Metadata struct {
	Kind     metadata.Element `json:"kind"`
	Value    string           `json:"value"`
	Ontology string           `json:"ontology"`
}

type Node struct {
	Body  string `json:"body"`
	Color string `json:"color"`
}

type Relation struct {
	From int               `json:"from"`
	To   int               `json:"to"`
	Kind relations.Element `json:"kind"`
}
