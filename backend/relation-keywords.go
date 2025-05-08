package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func RelationKeywords(c echo.Context) error {
	response := RelationKeywordsResponse{
		[]string{
			"test",
			"123",
		},
	}

	return c.JSON(http.StatusOK, response)
}
