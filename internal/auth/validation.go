// Pacakge auth - validation.go - Describe validation rules for authentication
package auth

import (
	"net/mail"
	"unicode"
)

func validateEmail(str string) error {
	if _, err := mail.ParseAddress(str); err != nil {
		return ErrInvalidEmail
	}
	return nil
}

func validatePassword(password string, confirm string) error {
	if password != confirm {
		return ErrPasswordMismatch
	}

	if len(password) < 12 {
		return ErrPasswordTooShort
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasDigit   bool
		hasSpecial bool
	)

	for _, r := range password {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasDigit = true
		case unicode.IsPunct(r), unicode.IsSymbol(r):
			hasSpecial = true
		}
	}

	switch {
	case !hasUpper:
		return ErrPasswordNoUppercase
	case !hasLower:
		return ErrPasswordNoLowercase
	case !hasDigit:
		return ErrPasswordNoDigit
	case !hasSpecial:
		return ErrPasswordNoSpecial
	}

	return nil
}
