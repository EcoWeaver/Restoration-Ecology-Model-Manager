package main

import "backend/relations"

type Graph struct {
	Nodes     []Node
	Relations []Relation
}

type Node struct {
	Body  string
	Color string
}

type Relation struct {
	From int
	To   int
	Kind relations.Element
}
