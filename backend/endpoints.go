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
		subCategory := element.SubCategory().String()
		category := element.Category().String()

		if response[category] == nil {
			response[category] = map[string][]string{}
		}

		if response[category][subCategory] == nil {
			response[category][subCategory] = []string{}
		}

		response[category][subCategory] = append(response[category][subCategory], element.String())
	}

	return c.JSON(http.StatusOK, response)
}

func GetRelationElements(c echo.Context) error {
	response := map[string][]string{}

	for _, element := range relations.ElementValues() {
		category := element.Category().String()

		if response[category] == nil {
			response[category] = []string{}
		}

		response[category] = append(response[category], element.String())
	}

	return c.JSON(http.StatusOK, response)
}
