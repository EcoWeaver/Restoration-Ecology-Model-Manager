package metadata

type Category int

//go:generate go run github.com/dmarkham/enumer -json -transform=kebab -type=Category
const (
	Core Category = iota
	Extended
)
