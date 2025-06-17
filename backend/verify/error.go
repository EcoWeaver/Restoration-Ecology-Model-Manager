package verify

var REQ = true

type Error struct{ E error }

func (self Error) Error() string {
	return self.E.Error()
}

func (self Error) Unwrap() error {
	return self.E
}
