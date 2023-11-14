package deeplink

import (
	"github.com/go-faster/errors"
)

// ValidateDomain validate given domain (user) name
func ValidateDomain(domain string) error {
	return checkDomainSymbols(domain)
}

// checkDomainSymbols check that domain contains only a-z, A-Z, 0-9 and '_'
// symbols.
func checkDomainSymbols(domain string) error {
	switch {
	case domain == "":
		return errors.New("is empty")
	case len(domain) > 32:
		return errors.New("is too big")
	case !IsLatinLower(rune(domain[0])):
		return errors.New("must start with lower letter")
	case domain[len(domain)-1] == '_':
		return errors.New("must not end with '_'")
	}

	for i, r := range domain {
		switch {
		case !IsLatinLower(r) && !IsDigit(r) && r != '_':
		case i > 0 && domain[i] == '_' && domain[i] == domain[i-1]:
		default:
			continue
		}

		return errors.Errorf("unexpected %c at %d", r, i)
	}

	return nil
}

func IsDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

// IsLatinLower reports whether the rune is lower ASCII latin letter.
func IsLatinLower(r rune) bool {
	r |= 0x20
	return 'a' <= r && r <= 'z'
}
