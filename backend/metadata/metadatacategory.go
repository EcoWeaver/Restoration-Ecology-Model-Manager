package metadata

type MetadataCategory int

//go:generate go run github.com/dmarkham/enumer -json -transform=kebab -type=MetadataCategory
const (
	Core MetadataCategory = iota
	Extended
)
