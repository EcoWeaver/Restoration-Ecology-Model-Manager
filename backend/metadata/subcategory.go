package metadata

type SubCategory int

//go:generate go run github.com/dmarkham/enumer -json -transform=kebab -type=SubCategory
const (
	ModelIdentification SubCategory = iota
	ModelDescription
	Specifications
	DataRequirements
	ModelComponents
	ModelValidation
	Usage
	References
	ModelReproducability
	ModelInteroperability
	UncertaintyAndSensitivity
	VersionControl
	CommunityEngagement
)

func (self SubCategory) Category() Category {
	switch self {
	case ModelIdentification,
		ModelDescription,
		Specifications,
		DataRequirements,
		ModelComponents,
		ModelValidation,
		Usage,
		References:
		return Core
	case ModelReproducability,
		ModelInteroperability,
		UncertaintyAndSensitivity,
		VersionControl,
		CommunityEngagement:
		return Extended
	default:
		panic("unreachable")
	}
}
