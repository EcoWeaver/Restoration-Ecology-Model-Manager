package metadata

type Element int

//go:generate go run github.com/dmarkham/enumer -json -text -transform=kebab -type=Element
const (
	Creator Element = iota + 1
	Title
	Publisher
	PublicationYear
	Subject
	Contributor
	Date
	Language
	ResourceType
	AlternateIdentifier
	RelatedIdentifier
	Size
	Format
	Version
	Rights
	Description
	GeoLocation
	FundingReference
	RelatedItem
)

func (self Element) Obligation() Obligation {
	switch self {
	case Creator, Title, Publisher, PublicationYear, ResourceType:
		return Mandatory
	case Subject, Contributor, Date, RelatedIdentifier, Description, GeoLocation:
		return Recommended
	case Language, AlternateIdentifier, Size, Format, Version, Rights, FundingReference, RelatedItem:
		return Optional
	default:
		panic("metadata element " + self.String() + " doesn't have an assigned obligation")
	}
}
