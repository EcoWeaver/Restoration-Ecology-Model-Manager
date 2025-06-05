package verify

import (
	"errors"
	"net/http"
	"net/url"
	"regexp"
)

var orcidRegexp = regexp.MustCompile(`\d{4}-\d{4}-\d{4}-\d{4}`)

func Orcid(s string) error {
	if s == "" {
		return errors.New("no orcid provided")
	}

	if !orcidRegexp.MatchString(s) {
		return errors.New("not a valid orcid")
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
		return errors.New("not a registered orcid")
	}

	return nil
}
