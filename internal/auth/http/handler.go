// Package http: handler.go - describe HTTP controller for auth purpose
package http

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/meplos/locana-deck-builder/internal/auth"
)

type Handler struct {
	registerUC *auth.RegisterUseCase
	loginUC    *auth.LoginUseCase
}

func NewHandler(registerUC *auth.RegisterUseCase, loginUC *auth.LoginUseCase) *Handler {
	return &Handler{
		registerUC: registerUC,
		loginUC:    loginUC,
	}
}

func (h *Handler) Register(ctx *echo.Context) error {
	var body RegisterRequest
	if err := ctx.Bind(&body); err != nil {
		return handleError(ctx, ErrInvalidRequest)
	}
	response, err := h.registerUC.Execute(ctx.Request().Context(), auth.RegisterInput{
		Name:            body.Name,
		Email:           body.Email,
		Password:        body.Password,
		ConfirmPassword: body.Confirm,
	})
	if err != nil {
		return handleError(ctx, err)
	}

	ctx.SetCookie(&http.Cookie{
		Name:     "ldb-tkn",
		Value:    response.Token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   86400,
	})

	return ctx.JSON(http.StatusAccepted, nil)
}

func (h *Handler) Login(ctx *echo.Context) error {
	var body LoginRequest
	if err := ctx.Bind(&body); err != nil {
		return handleError(ctx, ErrInvalidRequest)
	}
	response, err := h.loginUC.Execute(ctx.Request().Context(), auth.LoginInput{
		Email:    body.Email,
		Password: body.Password,
	})
	if err != nil {
		return handleError(ctx, err)
	}

	ctx.SetCookie(&http.Cookie{
		Name:     "ldb-tkn",
		Value:    response.Token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   86400,
	})

	return ctx.JSON(http.StatusAccepted, nil)
}
