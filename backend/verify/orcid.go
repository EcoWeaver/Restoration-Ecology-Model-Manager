package verify

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

var orcidRegexp = regexp.MustCompile(`^\d{4}-\d{4}-\d{4}-\d{3}[\dX]$`)
var orcidAccessToken *string

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

	if orcidAccessToken == nil {
		type AccessToken struct {
			AccessToken string `json:"access_token"`
		}

		reqBody := make(url.Values)
		reqBody.Add("client_id", os.Getenv("ORCID_CLIENT_ID"))
		reqBody.Add("client_secret", os.Getenv("ORCID_CLIENT_SECRET"))
		reqBody.Add("grant_type", os.Getenv("cli"))

		api := "https://pub.orcid.org/oauth/token"
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, api, strings.NewReader(reqBody.Encode()))
		if err != nil {
			return err
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		var t AccessToken
		if err := json.Unmarshal(body, &t); err != nil {
			return err
		}

		at := fmt.Sprintf("Bearer %s", t.AccessToken)
		orcidAccessToken = &at
	}

	api := fmt.Sprintf("https://pub.orcid.org/v2.1/%s/record", s)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, api, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", *orcidAccessToken)

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
