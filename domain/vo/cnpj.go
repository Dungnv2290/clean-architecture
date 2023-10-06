package vo

import (
	"errors"
	"regexp"
)

var (
	ErrInvalidCNPJ = errors.New("invalid cnpj")

	rxCNPJ = regexp.MustCompile(`^\d{2}\.?\d{3}\.?\d{3}\/?(:?\d{3}[1-9]|\d{2}[1-9]\d|\d[1-9]\d{2}|[1-9]\d{3})-?\d{2}$`)
)

// Cnpj structure
type Cnpj struct {
	value string
}

// NewCNPJ create new Cnpj
func NewCNPJ(value string) (Cnpj, error) {
	var c = Cnpj{value: value}

	if !c.validate() {
		return Cnpj{}, ErrInvalidCNPJ
	}

	return c, nil
}

func (c Cnpj) validate() bool {
	return rxCNPJ.MatchString(c.value)
}

// Value return value of Cnpj
func (c Cnpj) Value() string {
	return c.value
}

// String returns string representation of the Cnpj
func (c Cnpj) String() string {
	return c.value
}

// Equals check that two Cnpj are equal
func (c Cnpj) Equals(value Value) bool {
	o, ok := value.(Cnpj)
	return ok && c.value == o.value
}
