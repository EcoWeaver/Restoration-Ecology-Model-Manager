package verify

import (
	"context"
	"errors"
	"fmt"
	"regexp"
)

var orcidRegexp = regexp.MustCompile(`^\d{4}-\d{4}-\d{4}-\d{3}[\dX]$`)

func ORCID(s string, ctx context.Context) error {
	if s == "" {
		return Error{E: errors.New("no orcid provided")}
	}

	if !orcidRegexp.MatchString(s) {
		return Error{E: fmt.Errorf("encountered invalid orcid '%s'", s)}
	}

	return nil
}
