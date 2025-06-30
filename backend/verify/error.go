package verify

import (
	"fmt"
)

var REQ = true

func errNotProvided(t string) Error {
	return Error{E: fmt.Errorf("%s not provided", t)}
}

func errInvalid(t string) Error {
	return Error{E: fmt.Errorf("invalid %s", t)}
}

func errUnregistered(t string) Error {
	return Error{E: fmt.Errorf("unregistered %s", t)}
}

type Error struct{ E error }

func (self Error) Error() string {
	return self.E.Error()
}

func (self Error) Unwrap() error {
	return self.E
}
