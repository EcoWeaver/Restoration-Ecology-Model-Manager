package verify

import "regexp"

var colorRegexp = regexp.MustCompile(`(?i)^#[A-F0-9]{6}$`)

func Color(s string) error {
	if s == "" {
		return errNotProvided("color")
	}

	if !colorRegexp.MatchString(s) {
		return errInvalid("color")
	}

	return nil
}
