package verify

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
)

var doiRegexp = regexp.MustCompile(`(?i)^10\.\d{4,9}/[-._;()/:A-Z0-9]+$`)

func doi(s string) error {
	if s == "" {
		return Error{E: errors.New("no doi provided")}
	}

	if !doiRegexp.MatchString(s) {
		return Error{E: fmt.Errorf("encountered invalid doi '%s'", s)}
	}

	if !REQ {
		return nil
	}

	res, err := http.Get("https://doi.org/" + url.QueryEscape(s))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return Error{E: fmt.Errorf("encountered unregistered doi '%s'", s)}
	}

	return nil
}
