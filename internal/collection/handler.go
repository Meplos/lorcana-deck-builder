package collection

import (
	"context"
	"io"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/meplos/locana-deck-builder/internal/collection/domain"
	"github.com/meplos/locana-deck-builder/internal/db/mongo/schema"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Handler struct {
	ctx       context.Context
	listUC    *domain.ListUseCase
	createUC  *domain.CreateUseCase
	addCardUC *domain.AddCardUseCase
	exportUC  *domain.ExportUseCase
}

func Init(ctx context.Context, e *echo.Group, listUC *domain.ListUseCase, createUC *domain.CreateUseCase, addCardUC *domain.AddCardUseCase, exportUC *domain.ExportUseCase) {
	h := Handler{
		ctx:       ctx,
		listUC:    listUC,
		createUC:  createUC,
		addCardUC: addCardUC,
		exportUC:  exportUC,
	}

	e.GET("/collections", h.List)
	e.POST("/collections", h.Create)
	e.POST("/collections/card", h.AddCard)
	e.GET("/collections/export", h.Export)
}

type CollectionResponse struct {
	Name  string
	Cards schema.OwnedCard
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

	res := h.listUC.List(domain.PaginateParams{
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

	err := h.createUC.Create(body.Name)
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

	if err := h.addCardUC.AddCardToCollection(collID, body.CardID, body.Quantity); err != nil {
		return ctx.JSON(http.StatusInternalServerError, nil)
	}

	return ctx.JSON(http.StatusCreated, nil)
}

func (h *Handler) Export(ctx *echo.Context) error {
	file, err := h.exportUC.ExportAll()
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
