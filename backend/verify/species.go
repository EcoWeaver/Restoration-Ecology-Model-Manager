package verify

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

func species(s string, ctx context.Context) error {
	if s == "" {
		return errNotProvided("species")
	}

	if !REQ {
		return nil
	}

	api := "https://api.gbif.org/v2/species/match?scientificName=" + url.QueryEscape(s)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, api, nil)
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

	type Response struct {
		Usage struct {
			Status string `json:"status"`
		} `json:"usage"`
		Diagnostics struct {
			MatchType string `json:"matchType"`
		} `json:"diagnostics"`
	}

	var r Response
	if err = json.Unmarshal(body, &r); err != nil {
		return err
	}

	if r.Usage.Status != "ACCEPTED" || r.Diagnostics.MatchType != "EXACT" {
		return errUnregistered("species")
	}

	return nil
}
