package relations

type Category int

//go:generate go run github.com/dmarkham/enumer -json -text -transform=kebab -type=Category
const (
	Causal Category = iota
	Noncausal
)
