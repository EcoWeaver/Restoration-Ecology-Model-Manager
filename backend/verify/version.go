package verify

import (
	"errors"
	"fmt"
	"regexp"
)

var versionRegexp = regexp.MustCompile(`^1(\.\d+)*$`)

func Version(s string) error {
	if s == "" {
		return Error{E: errors.New("no version provided")}
	}

	if !versionRegexp.MatchString(s) {
		return Error{E: fmt.Errorf("encountered invalid version '%s'", s)}
	}

	return nil
}
