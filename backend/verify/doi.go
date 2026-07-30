package verify

import (
	"context"
	"net/http"
	"net/url"
	"regexp"
)

var doiRegexp = regexp.MustCompile(`(?i)^10\.\d{4,9}/[-._;()/:A-Z0-9]+$`)

func doi(s string, ctx context.Context, noVerifyReq bool) error {
	if s == "" {
		return errNotProvided("doi")
	}

	if !doiRegexp.MatchString(s) {
		return errInvalid("doi")
	}

	if noVerifyReq {
		return nil
	}

	api := "https://doi.org/api/handles/" + url.PathEscape(s)
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
		return errUnregistered("doi")
	}

	return nil
}
