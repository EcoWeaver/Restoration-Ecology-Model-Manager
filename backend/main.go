package main

import (
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
)

var (
	rw    = sync.RWMutex{}
	state = "Hello, World!"
)

func main() {
	e := echo.New()

	e.GET("/set", set)
	e.GET("/get", get)

	e.Logger.Fatal(e.Start(":8080"))
}

func set(c echo.Context) error {
	set := c.QueryString()

	rw.Lock()
	state = set
	rw.Unlock()

	return c.String(http.StatusOK, set)
}

func get(c echo.Context) error {
	rw.RLock()
	get := state
	rw.RUnlock()

	return c.String(http.StatusOK, get)
}
