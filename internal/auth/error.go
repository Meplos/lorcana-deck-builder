// Pacakge auth - error.go - Describe domain specific error
package auth

import "errors"

var (
	ErrInvalidEmail        = errors.New("invalid email")
	ErrPasswordTooShort    = errors.New("password must contain at least 12 characters")
	ErrPasswordMismatch    = errors.New("passwords do not match")
	ErrPasswordNoUppercase = errors.New("password must contain an uppercase letter")
	ErrPasswordNoLowercase = errors.New("password must contain a lowercase letter")
	ErrPasswordNoDigit     = errors.New("password must contain a digit")

	ErrPasswordNoSpecial = errors.New("password must contain a special character")
	ErrAlreadyExist      = errors.New("user already exist")
)
