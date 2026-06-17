package http

import (
	"io"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/meplos/locana-deck-builder/internal/collection"
	"github.com/meplos/locana-deck-builder/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Handler struct {
	listUC    *collection.ListUseCase
	createUC  *collection.CreateUseCase
	addCardUC *collection.AddCardUseCase
	exportUC  *collection.ExportUseCase
}

func NewHandler(
	listUC *collection.ListUseCase,
	createUC *collection.CreateUseCase,
	addCardUC *collection.AddCardUseCase,
	exportUC *collection.ExportUseCase,
) *Handler {
	return &Handler{
		listUC:    listUC,
		createUC:  createUC,
		addCardUC: addCardUC,
		exportUC:  exportUC,
	}
}

type CollectionResponse struct {
	Name  string
	Cards []domain.CollectionCard
}

type PaginateRequest struct {
	Page  int    `query:"page"`
	Limit int    `query:"limit"`
	Name  string `query:"name"`
}

type PaginateResponse struct {
	Page  int                  `json:"page"`
	Total int                  `json:"total"`
	Docs  []OwnedCollectionDTO `json:"docs"`
	Size  int                  `json:"size"`
}

func (h *Handler) List(ctx *echo.Context) error {
	query := PaginateRequest{}

	if err := ctx.Bind(&query); err != nil {
		return ctx.JSON(http.StatusBadRequest, nil)
	}

	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 20
	}

	res := h.listUC.List(ctx.Request().Context(), collection.PaginateParams{
		Limit: query.Limit,
		Page:  query.Page,
		Name:  query.Name,
	})

	ownedDto := make([]OwnedCollectionDTO, 0)
	for _, c := range res.Docs {
		ownedDto = append(ownedDto, BuildOwnedCollectionDTO(c))
	}

	return ctx.JSON(http.StatusOK, PaginateResponse{
		Size:  res.Size,
		Total: res.Total,
		Page:  res.Page,
		Docs:  ownedDto,
	})
}

type CreateBodyRequest struct {
	Name string `json:"name"`
}

func (h *Handler) Create(ctx *echo.Context) error {
	body := CreateBodyRequest{}
	if err := ctx.Bind(&body); err != nil {
		return ctx.JSON(http.StatusBadRequest, nil)
	}

	err := h.createUC.Create(ctx.Request().Context(), body.Name)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, nil)
	}

	return ctx.JSON(http.StatusCreated, nil)
}

type AddCardBody struct {
	CollectionID string `json:"collectionId"`
	CardID       string `json:"cardId"`
	Quantity     int    `json:"quantity"`
}

func (h *Handler) AddCard(ctx *echo.Context) error {
	body := AddCardBody{}
	if err := ctx.Bind(&body); err != nil {
		return ctx.JSON(http.StatusBadRequest, nil)
	}

	collID, err := primitive.ObjectIDFromHex(body.CollectionID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, nil)
	}

	if err := h.addCardUC.AddCardToCollection(ctx.Request().Context(), collID, body.CardID, body.Quantity); err != nil {
		return ctx.JSON(http.StatusInternalServerError, nil)
	}

	return ctx.JSON(http.StatusCreated, nil)
}

func (h *Handler) Export(ctx *echo.Context) error {
	file, err := h.exportUC.ExportAll(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, nil)
	}
	defer file.Close()
	ctx.Response().Header().Set(
		echo.HeaderContentDisposition,
		`attachment; filename="export.csv"`,
	)
	ctx.Response().Header().Set(echo.HeaderContentType, "text/csv; charset=utf-8")

	_, err = io.Copy(ctx.Response(), file)
	return err
}
