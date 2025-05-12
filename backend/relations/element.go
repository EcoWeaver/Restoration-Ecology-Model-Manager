package relations

type Element int

//go:generate go run github.com/dmarkham/enumer -json -transform=kebab -type=Element
const (
	CorrelatedWith Element = iota
	Increases
	Decreases
)

func (self Element) Category() Category {
	switch self {
	case CorrelatedWith:
		return Causal
	case Increases, Decreases:
		return Noncausal
	default:
		panic("unreachable")
	}
}
