package verify

import (
	"errors"
	"fmt"
	"regexp"
)

var colorRegexp = regexp.MustCompile(`(?i)^#[A-F0-9]{6}$`)

func Color(s string) error {
	if s == "" {
		return Error{E: errors.New("no color provided")}
	}

	if !colorRegexp.MatchString(s) {
		return Error{E: fmt.Errorf("encountered invalid color '%s'", s)}
	}

	return nil
}
