package deck

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/meplos/locana-deck-builder/internal/deck/domain"
	"github.com/meplos/locana-deck-builder/internal/ink"
)

type Handler struct {
	buildUC *domain.BuildDeckUseCase
}

func Init(e *echo.Group, buildUC *domain.BuildDeckUseCase) {
	h := Handler{
		buildUC: buildUC,
	}

	e.POST("/deck", h.BuildDeck)
}

type BuildDeckBody struct {
	Collection string         `json:"collection"`
	Colors     []ink.InkColor `json:"colors"`
	Level      string         `json:"level"`
}

func (h *Handler) BuildDeck(ctx *echo.Context) error {
	var body BuildDeckBody

	if err := ctx.Bind(&body); err != nil {
		return ctx.JSON(http.StatusBadRequest, nil)
	}

	res, err := h.buildUC.Build(domain.BuildDeckParams{
		CollectionName: body.Collection,
		InkColors:      body.Colors,
		Level:          body.Level,
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, res)
}
