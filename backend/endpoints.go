package main

import (
	"backend/metadata"
	"backend/relations"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetMetadataElements(c echo.Context) error {
	response := map[string]map[string][]string{}

	for _, element := range metadata.ElementValues() {
		category := element.Category().String()
		subCategory := element.SubCategory().String()

		if response[category] == nil {
			response[category] = map[string][]string{}
		}

		response[category][subCategory] = append(response[category][subCategory], element.String())
	}

	return c.JSON(http.StatusOK, response)
}

func GetRelationElements(c echo.Context) error {
	response := map[string][]string{}

	for _, element := range relations.ElementValues() {
		category := element.Category().String()
		response[category] = append(response[category], element.String())
	}

	return c.JSON(http.StatusOK, response)
}
