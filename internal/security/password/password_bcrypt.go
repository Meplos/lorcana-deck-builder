// Package password - password_bcrypt.go - Implement password utilities with bcrypt
package password

import "golang.org/x/crypto/bcrypt"

type BcryptHasher struct{}

func (h *BcryptHasher) Hash(str []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(str, bcrypt.DefaultCost)
}

func (h *BcryptHasher) Compare(hash []byte, password []byte) bool {
	if err := bcrypt.CompareHashAndPassword(hash, password); err != nil {
		return false
	}
	return true
}
