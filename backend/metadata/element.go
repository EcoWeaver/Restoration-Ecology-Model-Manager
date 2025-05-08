package metadata

type Element int

//go:generate go run github.com/dmarkham/enumer -json -transform=kebab -type=Element
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
