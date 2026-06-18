package http

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/meplos/locana-deck-builder/internal/deck"
	"github.com/meplos/locana-deck-builder/internal/domain"
)

type Handler struct {
	buildUC *deck.BuildDeckUseCase
	saveUC  *deck.SaveUseCase
}

func NewHandler(buildUC *deck.BuildDeckUseCase, saveUC *deck.SaveUseCase) *Handler {
	return &Handler{
		buildUC: buildUC,
		saveUC:  saveUC,
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

	return ctx.JSON(http.StatusOK, DeckBuildResponse{
		Size:     res.Size,
		Strategy: res.Strategy,
		Deck:     MapDeckCards(res.Deck),
		Name:     res.Name,
	})
}

func (h *Handler) Save(ctx *echo.Context) error {
	req := new(SaveDeckRequest)
	if err := ctx.Bind(req); err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	cards := make([]deck.CardInfo, 0)
	for _, c := range req.Deck {
		cards = append(cards, deck.CardInfo(c))
	}

	if err := h.saveUC.Execute(ctx.Request().Context(), &deck.SaveDeckInput{
		Size:     req.Size,
		Name:     req.Name,
		Strategy: req.Strategy,
		Deck:     cards,
	}); err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.NoContent(http.StatusCreated)
}
