package metadata

import (
	"fmt"
)

type Element int

//go:generate go run github.com/dmarkham/enumer -json -transform=kebab -type=Element
const (
	Identifier Element = iota + 1
	Creator
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

func (self Element) Obligation() (Obligation, error) {
	switch self {
	case Identifier, Creator, Title, Publisher, PublicationYear, ResourceType:
		return Mandatory, nil
	case Subject, Contributor, Date, RelatedIdentifier, Description, GeoLocation:
		return Recommended, nil
	case Language, AlternateIdentifier, Size, Format, Version, Rights, FundingReference, RelatedItem:
		return Optional, nil
	default:
		return 0, fmt.Errorf("unknown metadata element %d encountered", self)
	}
}
