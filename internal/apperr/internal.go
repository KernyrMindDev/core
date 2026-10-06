package apperr

import "fmt"

type UnregisteredCreator struct {
	Type string
}

func (u UnregisteredCreator) Error() string {
	return fmt.Sprintf("%v unregistered in factory", u.Type)
}
