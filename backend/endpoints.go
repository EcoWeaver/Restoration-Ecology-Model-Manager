package main

import (
	"backend/metadata"
	"backend/relations"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetMetadataElements(c echo.Context) error {
	response := map[string][]string{}

	for _, element := range metadata.ElementValues() {
		obligation, err := element.Obligation()
		if err != nil {
			return err
		}
		obligationS := obligation.String()
		response[obligationS] = append(response[obligationS], element.String())
	}

	return c.JSON(http.StatusOK, response)
}

func GetRelationElements(c echo.Context) error {
	response := map[string][]string{}

	for _, element := range relations.ElementValues() {
		category, err := element.Category()
		if err != nil {
			return err
		}
		elementS := category.String()
		response[elementS] = append(response[elementS], element.String())
	}

	return c.JSON(http.StatusOK, response)
}
