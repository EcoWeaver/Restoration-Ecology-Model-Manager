package verify

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"regexp"

	"golang.org/x/oauth2/clientcredentials"
)

var orcidRegexp = regexp.MustCompile(`^\d{4}-\d{4}-\d{4}-\d{3}[\dX]$`)
var orcidOAuth2Config *clientcredentials.Config

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

	if orcidOAuth2Config == nil {
		orcidOAuth2Config = &clientcredentials.Config{
			ClientID:     os.Getenv("ORCID_CLIENT_ID"),
			ClientSecret: os.Getenv("ORCID_CLIENT_SECRET"),
			TokenURL:     "https://pub.orcid.org/oauth/token",
			Scopes:       []string{"/read-public"},
		}
	}

	api := fmt.Sprintf("https://pub.orcid.org/v2.1/%s/record", s)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, api, nil)
	if err != nil {
		return err
	}

	res, err := orcidOAuth2Config.Client(ctx).Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errUnregistered("orcid")
	}

	return nil
}
