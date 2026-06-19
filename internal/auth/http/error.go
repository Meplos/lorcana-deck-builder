package http

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/meplos/locana-deck-builder/internal/auth"
)

var ErrInvalidRequest = errors.New("malformed request body")

type ErrorResponse struct {
	Error string `json:"error"`
}

func handleError(ctx *echo.Context, err error) error {
	status, message := mapAuthError(err)

	if status >= http.StatusInternalServerError {
		ctx.Logger().Error("auth handler error", "err", err)
	} else {
		ctx.Logger().Debug("auth handler error", "err", err)
	}

	return ctx.JSON(status, ErrorResponse{Error: message})
}

func mapAuthError(err error) (int, string) {
	switch {
	case errors.Is(err, ErrInvalidRequest):
		return http.StatusBadRequest, "Requête invalide."

	case errors.Is(err, auth.ErrInvalidEmail):
		return http.StatusBadRequest, "Email invalide."

	case errors.Is(err, auth.ErrPasswordTooShort):
		return http.StatusBadRequest, "Le mot de passe doit contenir au moins 12 caractères."

	case errors.Is(err, auth.ErrPasswordMismatch):
		return http.StatusBadRequest, "Les mots de passe ne correspondent pas."

	case errors.Is(err, auth.ErrPasswordNoUppercase):
		return http.StatusBadRequest, "Le mot de passe doit contenir une majuscule."

	case errors.Is(err, auth.ErrPasswordNoLowercase):
		return http.StatusBadRequest, "Le mot de passe doit contenir une minuscule."

	case errors.Is(err, auth.ErrPasswordNoDigit):
		return http.StatusBadRequest, "Le mot de passe doit contenir un chiffre."

	case errors.Is(err, auth.ErrPasswordNoSpecial):
		return http.StatusBadRequest, "Le mot de passe doit contenir un caractère spécial."

	case errors.Is(err, auth.ErrAlreadyExist):
		return http.StatusConflict, "Un compte existe déjà avec cet email."

	case errors.Is(err, auth.ErrNotExist), errors.Is(err, auth.ErrWrongPassword):
		return http.StatusUnauthorized, "Email ou mot de passe incorrect."

	default:
		return http.StatusInternalServerError, "Une erreur interne est survenue."
	}
}
