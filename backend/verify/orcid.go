package verify

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
)

var orcidRegexp = regexp.MustCompile(`^\d{4}-\d{4}-\d{4}-\d{3}[\dX]$`)

func orcid(s string, ctx context.Context) error {
	if s == "" {
		return Error{E: errors.New("no orcid provided")}
	}

	if !orcidRegexp.MatchString(s) {
		return Error{E: fmt.Errorf("encountered invalid orcid '%s'", s)}
	}

	if !REQ {
		return nil
	}

	api := "https://orcid.org/" + url.QueryEscape(s)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, api, nil)
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return Error{E: fmt.Errorf("encountered unregistered orcid '%s'", s)}
	}

	return nil
}
