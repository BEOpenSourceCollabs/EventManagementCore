package utils

import (
	"net/mail"
	"regexp"
)

type Validatable interface {
	Validate() []string
}

// StringLengthInBounds returns true of the provided string 's' is greater than or equal to min and less than or equal to max.
func StringLengthInBounds(s string, min, max int) bool {
	return len(s) >= min && len(s) <= max
}

// IsEmail returns true if the provided 'email' contains only the address part of a valid RFC 5322 address.
func IsEmail(email string) bool {
	addr, err := mail.ParseAddress(email)
	return err == nil && addr.Address == email
}

// ContainsAlphabeticCharacters returns true if the provided string 's' alphabetic numbers.
func ContainsAlphabeticCharacters(s string) bool {
	return regexp.MustCompile(`[a-zA-Z]`).MatchString(s)
}

// ContainsNumbericCharacters returns true if the provided string 's' contains numbers.
func ContainsNumbericCharacters(s string) bool {
	return regexp.MustCompile(`[0-9]`).MatchString(s)
}

// ContainsNoneAlphabeticCharacters returns true if the provided string 's' contains anything other than alphabetic characters.
func ContainsNoneAlphabeticCharacters(s string) bool {
	return len(regexp.MustCompile(`[a-zA-Z]+`).FindString(s)) != len(s)
}

// ContainsWhitespace returns true if the provided string 's' contains any whitespace characters.
func ContainsWhitespace(s string) bool {
	return regexp.MustCompile(`[\s]`).MatchString(s)
}

// IsAlphaNumeric returns true if the provided string 's' contains both alphabetic and numberal characters.
func IsAlphaNumeric(s string) bool {
	return len(regexp.MustCompile(`[a-zA-Z0-9]+`).FindString(s)) == len(s) && len(s) > 0
}
