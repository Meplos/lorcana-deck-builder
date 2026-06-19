// Package http: request.go - describe Http request data contract
package http

type (
	RegisterRequest struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Confirm  string `json:"confirm_password"`
	}
	LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
)
