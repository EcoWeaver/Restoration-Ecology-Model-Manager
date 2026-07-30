package verify

import (
	"backend/types"
	"fmt"
)

func errNotProvided(t string) types.Error {
	return types.Error{E: fmt.Errorf("%s not provided", t)}
}

func errInvalid(t string) types.Error {
	return types.Error{E: fmt.Errorf("invalid %s", t)}
}

func errUnregistered(t string) types.Error {
	return types.Error{E: fmt.Errorf("unregistered %s", t)}
}

func formatOnto(self types.Onto) string {
	return fmt.Sprintf("subclass of '%s' in ontology '%s'", self.Label, self.Onto)
}
