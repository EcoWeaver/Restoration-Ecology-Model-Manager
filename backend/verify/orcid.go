package verify

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
)

var orcidRegexp = regexp.MustCompile(`^\d{4}-\d{4}-\d{4}-\d{3}[\dX]$`)

func orcid(s string) error {
	if s == "" {
		return Error{E: errors.New("no orcid provided")}
	}

	if !orcidRegexp.MatchString(s) {
		return Error{E: fmt.Errorf("encountered invalid orcid '%s'", s)}
	}

	if !REQ {
		return nil
	}

	res, err := http.Get("https://orcid.org/" + url.QueryEscape(s))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return Error{E: fmt.Errorf("encountered unregistered orcid '%s'", s)}
	}

	return nil
}
