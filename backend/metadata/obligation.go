package metadata

type Obligation int

//go:generate go run github.com/dmarkham/enumer -json -text -transform=kebab -type=Obligation
const (
	Mandatory Obligation = iota + 1
	Recommended
	Optional
)
