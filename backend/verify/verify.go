package verify

var REQ = true

type VerifyError struct{ E error }

func (self VerifyError) Error() string {
	return self.E.Error()
}

func (self VerifyError) Unwrap() error {
	return self.E
}
