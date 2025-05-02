package main

type MetadataCategory struct {
	MetadataCategoryId int
	MainCategory       string
}

type SubCategory struct {
	SubCategoryId int
	SubCategory   string
}

type Element struct {
	ElementId          int
	MetadataCategoryId int
	SubCategoryId      int
	ElementName        string
	ElementDescription string
}

type MetadataValues struct {
	ValueId    int
	ModelId    int
	ElementId  int
	FieldValue string
	Ontology   string
}
