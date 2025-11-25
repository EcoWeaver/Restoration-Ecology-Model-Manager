package types

type Model struct {
	Graph    Graph      `json:"graph"`
	Parent   *string    `json:"parent"`
	Metadata []Metadata `json:"metadata"`
}

type Graph struct {
	Nodes     []Node     `json:"nodes"`
	Relations []Relation `json:"relations"`
}

type Metadata struct {
	Kind  string `json:"kind"`
	Value string `json:"value"`
}

type Node struct {
	Body     string    `json:"body"`
	Color    string    `json:"color"`
	Position Cartesian `json:"position"`
}

type Relation struct {
	Kind     string `json:"kind"`
	From     int    `json:"from"`
	To       int    `json:"to"`
	Position Polar  `json:"position"`
}

type Cartesian struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Polar struct {
	R   float64 `json:"r"`
	Phi float64 `json:"phi"`
}

type OLS4 struct {
	Onto  string `json:"onto"`
	IRI   string `json:"iri"`
	Label string `json:"label"`
}

type RelationWeight struct {
	From    int        `json:"from"`
	To      int        `json:"to"`
	Formula Polynomial `json:"formula"`
}

type Polynomial []Coefficient

type Coefficient struct {
	Factor   float64 `json:"factor"`
	Exponent float64 `json:"exponent"`
}

type Error struct{ E error }

func (self Error) Error() string {
	return self.E.Error()
}

func (self Error) Unwrap() error {
	return self.E
}
