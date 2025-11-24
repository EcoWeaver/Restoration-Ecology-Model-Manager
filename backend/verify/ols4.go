package verify

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
)

type OLS4 struct {
	Onto  string `json:"onto"`
	IRI   string `json:"iri"`
	Label string `json:"label"`
}

func (self OLS4) verify(s string, ctx context.Context) error {
	if s == "" {
		return errNotProvided(fmt.Sprintf("ontology '%s'", self.Onto))
	}

	if !REQ {
		return nil
	}

	api := "https://www.ebi.ac.uk/ols4/api/search?fieldList=label&queryFields=label&exact=true&ontology=" + self.Onto + "&childrenOf=" + self.IRI + "&q=" + url.QueryEscape(s)
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

	type Label struct {
		Label string `json:"label"`
	}

	type Response struct {
		Response struct {
			Docs []Label `json:"docs"`
		} `json:"response"`
	}

	var r Response
	if err = json.Unmarshal(body, &r); err != nil {
		return err
	}

	if !slices.ContainsFunc(r.Response.Docs, func(l Label) bool { return l.Label == s }) {
		return errInvalid(fmt.Sprintf("subclass of '%s' in ontology '%s'", self.Label, self.Onto))
	}

	return nil
}
