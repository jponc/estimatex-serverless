package ddbrepository

type ErrString string

func (e ErrString) Error() string {
	return string(e)
}

const ErrNotFound = ErrString("not found")
