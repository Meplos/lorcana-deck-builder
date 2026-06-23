package jwt

import "os"

type JWTManager interface {
	Create(ID string) (string, error)
	Parse(token string) (string, error)
}

func New() JWTManager {
	return &V5JWTManager{
		secret: os.Getenv("JWT_SECRET"),
	}
}
