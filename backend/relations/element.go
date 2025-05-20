package relations

type Element int

//go:generate go run github.com/dmarkham/enumer -json -text -transform=kebab -type=Element
const (
	CorrelatedWith Element = iota + 1
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
		panic("relations element " + self.String() + " doesn't have an assigned category")
	}
}
