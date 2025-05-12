package relations

type Category int

//go:generate go run github.com/dmarkham/enumer -json -transform=kebab -type=Category
const (
	Causal Category = iota
	Noncausal
)
