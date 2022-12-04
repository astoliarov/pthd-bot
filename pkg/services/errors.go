package services

import "fmt"

type ErrMixedCharactersInName struct {
	Name string
}

func NewErrMixedCharactersInName(name string) *ErrMixedCharactersInName {
	return &ErrMixedCharactersInName{
		Name: name,
	}
}

func (err ErrMixedCharactersInName) Error() string {
	return fmt.Sprintf("mixed characters in name: %s", err.Name)
}
