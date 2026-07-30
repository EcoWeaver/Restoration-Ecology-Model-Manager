package verify

import (
	"backend/util"
	"context"
	"net/url"
)

func species(s string, ctx context.Context, noVerifyReq bool) error {
	if s == "" {
		return errNotProvided("species")
	}

	if noVerifyReq {
		return nil
	}

	api := "https://api.gbif.org/v2/species/match?scientificName=" + url.QueryEscape(s)
	res, err := util.GetAndUnmarshal[struct {
		Usage struct {
			Status string `json:"status"`
		} `json:"usage"`
		Diagnostics struct {
			MatchType string `json:"matchType"`
		} `json:"diagnostics"`
	}](api, ctx)
	if err != nil {
		return err
	}

	if res.Usage.Status != "ACCEPTED" || res.Diagnostics.MatchType != "EXACT" {
		return errUnregistered("species")
	}

	return nil
}
