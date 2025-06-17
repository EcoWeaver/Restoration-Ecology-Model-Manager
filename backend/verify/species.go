package verify

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func species(s string) error {
	if s == "" {
		return Error{E: errors.New("no species provided")}
	}

	if !REQ {
		return nil
	}

	res, err := http.Get("https://api.gbif.org/v1/species/match?rank=SPECIES&name=" + url.QueryEscape(s))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	type Response struct {
		Status    string `json:"status"`
		MatchType string `json:"matchType"`
		Rank      string `json:"rank"`
	}

	var r Response
	if err = json.Unmarshal(body, &r); err != nil {
		return err
	}

	if r.Status != "ACCEPTED" || r.MatchType != "EXACT" || r.Rank != "SPECIES" {
		return Error{E: fmt.Errorf("encountered unregistered species '%s'", s)}
	}

	return nil
}
