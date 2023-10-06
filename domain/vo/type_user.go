package vo

import (
	"errors"
	"strings"
)

const (
	COMMON   TypeUser = "COMMON"
	MERCHANT TypeUser = "MERCHANT"
)

var (
	ErrInvalidTypeUser = errors.New("invalid type user")

	ErrNotAllowedTypeUser = errors.New("not allowed type user")
)

type (
	// TypeUser define user types
	TypeUser string
)

func NewTypeUser(value string) (TypeUser, error) {
	switch TypeUser(strings.ToUpper(value)) {
	case COMMON, MERCHANT:
		return TypeUser(strings.ToUpper(value)), nil
	}

	return "", ErrInvalidTypeUser
}

// String return string representation of the TypeUser
func (t TypeUser) String() string {
	return string(t)
}

// ToUpper
func (t TypeUser) ToUpper() TypeUser {
	return TypeUser(strings.ToUpper(string(t)))
}
