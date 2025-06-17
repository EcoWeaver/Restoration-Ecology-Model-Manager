package verify

import (
	"errors"
	"fmt"
	"regexp"
)

var verRegexp = regexp.MustCompile(`\d(\.\d+)*`)

func Version(s string) error {
	if s == "" {
		return Error{E: errors.New("no version provided")}
	}

	if !verRegexp.MatchString(s) {
		return Error{E: fmt.Errorf("encountered invalid version '%s'", s)}
	}

	return nil
}
