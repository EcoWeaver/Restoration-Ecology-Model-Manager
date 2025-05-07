package main

type MetadataCategory struct {
	MetadataCategoryId int
	MainCategory       string
}

type SubCategory struct {
	SubCategoryId      int
	MetadataCategoryId int
	SubCategory        string
}

type Element struct {
	ElementId          int
	SubCategoryId      int
	ElementName        string
	ElementDescription string
	IsMandatory        bool
	AllowedOntology    string
}

type MetadataValues struct {
	ValueId    int
	ModelId    int
	ElementId  int
	FieldValue string
	Ontology   string
}
