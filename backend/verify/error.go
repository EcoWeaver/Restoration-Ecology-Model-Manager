package verify

import (
	"backend/types"
	"fmt"
)

var REQ = true

func errNotProvided(t string) types.Error {
	return types.Error{E: fmt.Errorf("%s not provided", t)}
}

func errInvalid(t string) types.Error {
	return types.Error{E: fmt.Errorf("invalid %s", t)}
}

func errUnregistered(t string) types.Error {
	return types.Error{E: fmt.Errorf("unregistered %s", t)}
}
