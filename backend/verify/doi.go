package verify

import (
	"errors"
	"net/http"
	"net/url"
	"regexp"
)

var doiRegexp = regexp.MustCompile(`(?i)^10\.\d{4,9}/[-._;()/:A-Z0-9]+$`)

func Doi(s string) error {
	if s == "" {
		return errors.New("no doi provided")
	}

	if !doiRegexp.MatchString(s) {
		return errors.New("not a valid doi")
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
		return errors.New("not a registered doi")
	}

	return nil
}
