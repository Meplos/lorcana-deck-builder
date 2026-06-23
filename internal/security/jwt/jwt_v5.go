// Package jwt : jwt_v5.go - Implement JwtManager using jwt/v5
package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string `json:"uid"`
	jwt.RegisteredClaims
}

type V5JWTManager struct {
	secret string
}

func (j *V5JWTManager) Create(ID string) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		Claims{
			UserID: ID,
			RegisteredClaims: jwt.RegisteredClaims{
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			},
		},
	)

	signed, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return "", err
	}

	return signed, nil
}

func (j *V5JWTManager) Parse(token string) (string, error) {
	raw, err := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid algorithm")
		}
		return []byte(j.secret), nil
	})
	if err != nil {
		return "", err
	}
	claims, ok := raw.Claims.(*Claims)
	if !ok || !raw.Valid {
		return "", errors.New("invalid claims")
	}

	return claims.UserID, nil
}
