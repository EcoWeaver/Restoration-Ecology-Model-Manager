package verify

import (
	"errors"
	"fmt"
	"regexp"
)

var uuidRegexp = regexp.MustCompile(`(?i)^[A-F0-9]{4}([A-F0-9]{4}-){4}[A-F0-9]{12}$`)

func UUID(s string) error {
	if s == "" {
		return Error{E: errors.New("no uuid provided")}
	}

	if !uuidRegexp.MatchString(s) {
		return Error{E: fmt.Errorf("encountered invalid uuid '%s'", s)}
	}

	return nil
}
