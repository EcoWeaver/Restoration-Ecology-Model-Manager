package verify

import (
	"context"
	"net/http"
	"net/url"
	"regexp"
)

var orcidRegexp = regexp.MustCompile(`^\d{4}-\d{4}-\d{4}-\d{3}[\dX]$`)

func orcid(s string, ctx context.Context) error {
	if s == "" {
		return errNotProvided("orcid")
	}

	if !orcidRegexp.MatchString(s) {
		return errInvalid("orcid")
	}

	if !REQ {
		return nil
	}

	api := "https://pub.orcid.org/v3.0/" + url.PathEscape(s)
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
		return errUnregistered("orcid")
	}

	return nil
}
