package vo

import "errors"

const (
	// Currency types
	BRL TypeCurrency = "BRL"
	USD TypeCurrency = "USD"
)

type (
	// TypeCurrency define currency types
	TypeCurrency string
)

func (tc TypeCurrency) String() string {
	return string(tc)
}

var (
	ErrInvalidCurrency = errors.New("invalid currency")
)

type (
	// Currency structure
	Currency struct {
		value TypeCurrency
	}
)

// NewCurrency create new Currency
func NewCurrency(value string) (Currency, error) {
	var c = Currency{value: TypeCurrency(value)}

	if !c.validate() {
		return Currency{}, ErrInvalidCurrency
	}

	return c, nil
}

func (c Currency) validate() bool {
	switch c.value {
	case BRL, USD:
		return true
	}

	return false
}

// Value return value Currency
func (c Currency) Value() TypeCurrency {
	return c.value
}

// String return string representation of the Currency
func (c Currency) String() string {
	return string(c.value)
}

// Equals check that two Currency are the same
func (c Currency) Equals(value Value) bool {
	o, ok := value.(Currency)
	return ok && c.value == o.value
}
