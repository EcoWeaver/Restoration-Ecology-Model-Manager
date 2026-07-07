package verify

import (
	"backend/types"
	"backend/util"
	"context"
	"net/url"
	"slices"
)

func bioPortal(self types.Onto, s string, ctx context.Context, bioPortalAPIKey string, noVerifyReq bool) error {
	if s == "" {
		return errNotProvided(formatOnto(self))
	}

	if noVerifyReq {
		return nil
	}

	type PrefLabel struct {
		PrefLabel string `json:"prefLabel"`
	}

	api := "https://data.bioontology.org/search?require_exact=true&apikey=" + bioPortalAPIKey + "&ontology=" + self.Onto + "&subtree_root_id=" + self.IRI + "&q=" + url.QueryEscape(s)
	res, err := util.GetAndUnmarshal[struct {
		Collection []PrefLabel `json:"collection"`
	}](api, ctx)
	if err != nil {
		return err
	}

	if !slices.ContainsFunc(res.Collection, func(l PrefLabel) bool { return l.PrefLabel == s }) {
		return errInvalid(formatOnto(self))
	}

	return nil
}
