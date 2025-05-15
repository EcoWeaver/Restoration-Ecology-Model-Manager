package relations

import "fmt"

type Element int

//go:generate go run github.com/dmarkham/enumer -json -transform=kebab -type=Element
const (
	CorrelatedWith Element = iota + 1
	Increases
	Decreases
)

func (self Element) Category() (Category, error) {
	switch self {
	case CorrelatedWith:
		return Causal, nil
	case Increases, Decreases:
		return Noncausal, nil
	default:
		return 0, fmt.Errorf("unknown relations element %d encountered", self)
	}
}
