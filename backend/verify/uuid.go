package verify

import "regexp"

var uuidRegexp = regexp.MustCompile(`(?i)^[A-F0-9]{4}([A-F0-9]{4}-){4}[A-F0-9]{12}$`)

func UUID(s string) error {
	if s == "" {
		return errNotProvided("uuid")
	}

	if !uuidRegexp.MatchString(s) {
		return errInvalid("uuid")
	}

	return nil
}
