package verify

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
)

func Species(s *string) error {
	if *s == "" {
		return VerifyError{E: errors.New("no species provided")}
	}

	if !REQ {
		return nil
	}

	res, err := http.Get("https://api.gbif.org/v1/species/match?rank=SPECIES&name=" + url.QueryEscape(*s))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	type Response struct {
		Status        string `json:"status"`
		MatchType     string `json:"matchType"`
		Rank          string `json:"rank"`
		CanonicalName string `json:"canonicalName"`
	}

	var r Response
	if err = json.Unmarshal(body, &r); err != nil {
		return err
	}

	if r.Status != "ACCEPTED" && r.Status != "SYNONYM" && r.Status != "HOMOTYPIC_SYNONYM" {
		return VerifyError{E: errors.New("species name is not accepted")}
	}

	if r.MatchType != "EXACT" && r.MatchType != "FUZZY" {
		return VerifyError{E: errors.New("no match or multiple matches")}
	}

	if r.Rank != "SPECIES" {
		return VerifyError{E: errors.New("match is not a species")}
	}

	*s = r.CanonicalName

	return nil
}
