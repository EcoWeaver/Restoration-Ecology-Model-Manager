package main

type MetadataValues struct {
	ValueId    int
	ModelId    int
	ElementId  Element
	FieldValue string
	Ontology   string
}

type MetadataCategory int
type SubCategory int
type Element int

const (
	Core MetadataCategory = iota
	Extended
)

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

func (self SubCategory) MetadataCategory() MetadataCategory {
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

const (
	Uuid Element = iota
	ModelName
	Version
	CreationDate
	LastUpdateDate
	AuthorsAndContributors
	License
	Abstract
	Keywords
	ModelType
	EcologicalDomain
	SpatialScale
	TemporalScale
	TargetSpeciesOrEcosystem
	LanguageOrSoftware
	DependenciesAndLibraries
	HardwareRequirements
	OSCompatibility
	DataInputTypeFormat
	DataOutputTypeFormat
	DataSources
	DataPreprocessing
	Variables
	VariableRelationships
	Parameters
	EquationsAndAlgorithms
	CalibrationTechniques
	ValidationTechniques
	PerformanceMetrics
	SensitivityAnalysis
	InstallInstructions
	ManualOrDocs
	ExampleUseCases
	LimitationsOrAssumptions
	RelatedPublications
	Citations
	RepoLink
	DockerInfo
	ComputationalEnvironment
	ApiDocs
	DataExchangeFormats
	ModelAndToolIntegration
	UncertaintyQuantificationMethods
	Robustness
	Changelog
	PreviousVersions
	VersionCompatibility
	UserCommunity
	ForumsOrDiscussions
	ContributionGuidelines
)

func (self Element) SubCategory() SubCategory {
	switch self {
	case Uuid,
		ModelName,
		Version,
		CreationDate,
		LastUpdateDate,
		AuthorsAndContributors,
		License:
		return ModelIdentification
	case Abstract,
		Keywords,
		ModelType,
		EcologicalDomain,
		SpatialScale,
		TemporalScale,
		TargetSpeciesOrEcosystem:
		return ModelDescription
	case LanguageOrSoftware,
		DependenciesAndLibraries,
		HardwareRequirements,
		OSCompatibility:
		return Specifications
	case DataInputTypeFormat,
		DataOutputTypeFormat,
		DataSources,
		DataPreprocessing:
		return DataRequirements
	case Variables,
		VariableRelationships,
		Parameters,
		EquationsAndAlgorithms:
		return ModelComponents
	case CalibrationTechniques,
		ValidationTechniques,
		PerformanceMetrics,
		SensitivityAnalysis:
		return ModelValidation
	case InstallInstructions,
		ManualOrDocs,
		ExampleUseCases,
		LimitationsOrAssumptions:
		return Usage
	case RelatedPublications,
		Citations:
		return References
	case RepoLink,
		DockerInfo,
		ComputationalEnvironment:
		return ModelReproducability
	case ApiDocs,
		DataExchangeFormats,
		ModelAndToolIntegration:
		return ModelInteroperability
	case UncertaintyQuantificationMethods,
		Robustness:
		return UncertaintyAndSensitivity
	case Changelog,
		PreviousVersions,
		VersionCompatibility:
		return VersionControl
	case UserCommunity,
		ForumsOrDiscussions,
		ContributionGuidelines:
		return CommunityEngagement
	default:
		panic("unreachable")
	}
}

func (self Element) MetadataCategory() MetadataCategory {
	return self.SubCategory().MetadataCategory()
}
