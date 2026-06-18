// Package jwt : jwt_v5.go - Implement JwtManager using jwt/v5
package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type V5JWTManager struct {
	secret string
}

func (j *V5JWTManager) Create(ID string) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"uid": ID,
			"iat": time.Now().Unix(),
			"exp": time.Now().Add(time.Hour * 24),
		},
	)

	signed, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return "", err
	}

	return signed, nil
}
