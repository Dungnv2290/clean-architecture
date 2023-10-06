package vo

import (
	"errors"
	"regexp"
)

var (
	ErrInvalidEmail = errors.New("invalid email")

	rxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

// Email structure
type Email struct {
	value string
}

// NewEmail create new Email
func NewEmail(value string) (Email, error) {
	var e = Email{value: value}

	if !e.validate() {
		return Email{}, ErrInvalidEmail
	}

	return e, nil
}

func (e Email) validate() bool {
	return rxEmail.MatchString(e.value)
}

// Value return value Email
func (e Email) Value() string {
	return e.value
}

// String return string representation of the Email
func (e Email) String() string {
	return e.value
}

// Equals check that two Email are the same
func (e Email) Equals(value Value) bool {
	o, ok := value.(Email)
	return ok && e.value == o.value
}

// NewEmailTest create new Email for testing
func NewEmailTest(e string) Email {
	return Email{value: e}
}
