package cards

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/meplos/locana-deck-builder/internal/cards/domain"
)

type Handler struct {
	ctx    context.Context
	ListUC *domain.ListCardUC
}

func Init(e *echo.Group, listUc *domain.ListCardUC) {
	h := Handler{
		ListUC: listUc,
	}

	e.GET("/cards", h.ListCard)
}

func (h *Handler) ListCard(ctx *echo.Context) error {
	var body PaginateRequest
	if err := ctx.Bind(&body); err != nil {
		return ctx.JSON(http.StatusBadRequest, nil)
	}

	if body.PageNumber == 0 {
		body.PageNumber = 1
	}
	if body.Limit == 0 {
		body.Limit = 20
	}

	response := h.ListUC.ListCards(domain.ListCardParams{
		Page:   int64(body.PageNumber),
		Limit:  int64(body.Limit),
		Search: body.Search,
		Color:  body.Color,
	})

	cards := make([]CardDTO, 0)
	for _, d := range response.Docs {
		cards = append(cards, BuildCardDTO(d))
	}

	return ctx.JSON(http.StatusOK, PaginateResponse{
		Page:  body.PageNumber,
		Docs:  cards,
		Total: uint(response.Total),
		Size:  uint(response.Size),
	})
}
