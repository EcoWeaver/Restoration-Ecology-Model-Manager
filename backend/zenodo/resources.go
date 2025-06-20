package zenodo

import (
	"net/url"
	"time"

	"golang.org/x/text/language"
)

type Deposit struct {
	Created   time.Time         `json:"created"`
	Doi       *string           `json:"doi"`
	DoiUrl    *string           `json:"doi_url"`
	Files     []DepositionFile  `json:"files"`
	Id        int               `json:"id"`
	Metadata  []DepositMetadata `json:"metadata"`
	Modified  time.Time         `json:"modified"`
	Owner     int               `json:"owner"`
	RecordId  *int              `json:"record_id"`
	RecordUrl *url.URL          `json:"record_url"`
	State     string            `json:"state"`
	Submitted bool              `json:"submitted"`
	Title     string            `json:"title"`
}

type DepositMetadata struct {
	UploadType      string    `json:"upload_type"`
	PublicationType *string   `json:"publication_type"`
	ImageType       *string   `json:"image_type"`
	PublicationDate time.Time `json:"publication_date"`
	Title           string    `json:"title"`
	Creators        []struct {
		Name        string  `json:"name"`
		Affiliation *string `json:"affiliation"`
		Orcid       *string `json:"orcid"`
		Gnd         *string `json:"gnd"`
	} `json:"creators"`
	Description        string     `json:"description"`
	AccessRight        string     `json:"access_right"`
	License            *string    `json:"license"`
	EmbargoDate        *time.Time `json:"embargo_date"`
	AccessConditions   *string    `json:"access_conditions"`
	Doi                *string    `json:"doi"`
	PreserveDoi        *string    `json:"preserve_doi"`
	Keywords           *[]string  `json:"keywords"`
	Notes              *string    `json:"notes"`
	RelatedIdentifiers *[]struct {
		Identifier   string `json:"identifier"`
		Relation     string `json:"relation"`
		ResourceType string `json:"resource_type"`
	} `json:"related_identifiers"`
	Contributors *[]struct {
		Name        string  `json:"name"`
		Affiliation *string `json:"affiliation"`
		Orcid       *string `json:"orcid"`
		Gnd         *string `json:"gnd"`
	} `json:"contributors"`
	References  *[]string `json:"references"`
	Communities *[]struct {
		Identifier string `json:"identifier"`
	} `json:"communities"`
	Grants *[]struct {
		Id string `json:"id"`
	} `json:"grants"`
	JournalTitle          *string  `json:"journal_title"`
	JournalVolume         *string  `json:"journal_volume"`
	JournalIssue          *string  `json:"journal_issue"`
	JournalPages          *string  `json:"journal_pages"`
	ConferenceTitle       *string  `json:"conference_title"`
	ConferenceAcronym     *string  `json:"conference_acronym"`
	ConferenceDates       *string  `json:"conference_dates"`
	ConferencePlace       *string  `json:"conference_place"`
	ConferenceUrl         *url.URL `json:"conference_url"`
	ConferenceSession     *string  `json:"conference_session"`
	ConferenceSessionPart *string  `json:"conference_session_part"`
	ImprintPublisher      *string  `json:"imprint_publisher"`
	ImprintIsbn           *string  `json:"imprint_isbn"`
	ImprintPlace          *string  `json:"imprint_place"`
	PartofTitle           *string  `json:"partof_title"`
	PartofPages           *string  `json:"partof_pages"`
	ThesisSupervisors     *[]struct {
		Name        string  `json:"name"`
		Affiliation *string `json:"affiliation"`
		Orcid       *string `json:"orcid"`
		Gnd         *string `json:"gnd"`
	} `json:"thesis_supervisors"`
	ThesisUniversity *string `json:"thesis_university"`
	Subjects         *[]struct {
		Term       string `json:"term"`
		Identifier string `json:"identifier"`
		Scheme     string `json:"scheme"`
	} `json:"subjects"`
	Version   *string       `json:"version"`
	Language  *language.Tag `json:"Language"`
	Locations *[]struct {
		Lat         *float64 `json:"lat"`
		Lon         *float64 `json:"lon"`
		Place       *string  `json:"place"`
		Description *string  `json:"description"`
	} `json:"locations"`
	Dates *[]struct {
		Start       time.Time `json:"start"`
		End         time.Time `json:"end"`
		Type        string    `json:"type"`
		Description *string   `json:"description"`
	} `json:"dates"`
	Method *string `json:"method"`
}

type DepositionFile struct {
	Id       string `json:"id"`
	FileName string `json:"filename"`
	FileSize int    `json:"filesize"`
	Checksum string `json:"checksum"`
}
