package metadata

type Obligation int

//go:generate go run github.com/dmarkham/enumer -json -transform=kebab -type=Obligation
const (
	Mandatory Obligation = iota + 1
	Recommended
	Optional
)
