// Package password - password.go Handle security concern on Hashing Compare etc..
package password

type Hasher interface {
	Hash(str []byte) ([]byte, error)
	Compare(hash []byte, password []byte) bool
}

func New() Hasher {
	return &BcryptHasher{}
}
