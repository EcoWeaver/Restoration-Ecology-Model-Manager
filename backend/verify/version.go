package verify

import (
	"errors"
	"fmt"
	"regexp"
)

var verRegexp = regexp.MustCompile(`\d(\.\d+)*`)

func Version(v string) error {
	if v == "" {
		return errors.New("no version provided")
	}

	if !verRegexp.MatchString(v) {
		return fmt.Errorf("encountered invalid version '%s'", v)
	}

	return nil
}
