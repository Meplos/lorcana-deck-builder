package deck

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/meplos/locana-deck-builder/internal/deck"
	"github.com/meplos/locana-deck-builder/internal/domain"
)

type Handler struct {
	buildUC *deck.BuildDeckUseCase
}

func NewHandler(buildUC *deck.BuildDeckUseCase) *Handler {
	return &Handler{
		buildUC: buildUC,
	}
}

type BuildDeckBody struct {
	Collection string            `json:"collection"`
	Colors     []domain.InkColor `json:"colors"`
	Level      string            `json:"level"`
}

func (h *Handler) BuildDeck(ctx *echo.Context) error {
	var body BuildDeckBody

	if err := ctx.Bind(&body); err != nil {
		return ctx.JSON(http.StatusBadRequest, nil)
	}

	res, err := h.buildUC.Build(ctx.Request().Context(), deck.BuildDeckParams{
		CollectionName: body.Collection,
		InkColors:      body.Colors,
		Level:          body.Level,
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, res)
}
