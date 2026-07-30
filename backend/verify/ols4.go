package verify

import (
	"backend/types"
	"backend/util"
	"context"
	"net/url"
	"slices"
)

func ols4(self types.Onto, s string, ctx context.Context, noVerifyReq bool) error {
	if s == "" {
		return errNotProvided(formatOnto(self))
	}

	if noVerifyReq {
		return nil
	}

	type Label struct {
		Label string `json:"label"`
	}

	api := "https://www.ebi.ac.uk/ols4/api/search?fieldList=label&queryFields=label&exact=true&ontology=" + self.Onto + "&childrenOf=" + self.IRI + "&q=" + url.QueryEscape(s)
	res, err := util.GetAndUnmarshal[struct {
		Response struct {
			Docs []Label `json:"docs"`
		} `json:"response"`
	}](api, ctx)
	if err != nil {
		return err
	}

	if !slices.ContainsFunc(res.Response.Docs, func(l Label) bool { return l.Label == s }) {
		return errInvalid(formatOnto(self))
	}

	return nil
}
