package main

import (
	"backend/metadata"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetRelationKeywords(c echo.Context) error {
	return c.JSON(http.StatusOK, []string{
		"test",
		"123",
	})
}

func GetMetadataElements(c echo.Context) error {
	response := map[string]map[string][]string{}

	for _, element := range metadata.ElementValues() {
		subCategory := element.SubCategory().String()
		metadataCategory := element.MetadataCategory().String()

		if response[metadataCategory] == nil {
			response[metadataCategory] = map[string][]string{}
		}

		if response[metadataCategory][subCategory] == nil {
			response[metadataCategory][subCategory] = []string{}
		}

		response[metadataCategory][subCategory] = append(response[metadataCategory][subCategory], element.String())
	}

	return c.JSON(http.StatusOK, response)
}
