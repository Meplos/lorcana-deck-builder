package app

import (
	"github.com/meplos/locana-deck-builder/assets/images"
	"github.com/meplos/locana-deck-builder/internal/cards"
	cardsHttp "github.com/meplos/locana-deck-builder/internal/cards/http"
	"github.com/meplos/locana-deck-builder/internal/collection"
	collectionHttp "github.com/meplos/locana-deck-builder/internal/collection/http"
	"github.com/meplos/locana-deck-builder/internal/deck"
	deckHttp "github.com/meplos/locana-deck-builder/internal/deck/http"
	"go.mongodb.org/mongo-driver/mongo"
)

type Container struct {
	CardHandler       *cardsHttp.Handler
	CollectionHandler *collectionHttp.Handler
	DeckHandler       *deckHttp.Handler
}

func NewContainer(DB *mongo.Database) (*Container, error) {
	imageURIBuilder := images.New()

	cardsRepository := cards.SetupMongoCardRepository(DB)
	cardListUC := cards.CreateListCardUC(cardsRepository)

	cardHandler := cardsHttp.NewHandler(cardListUC)

	collectionRepo := collection.NewCollectionRepo(DB)
	collectionListUC := collection.NewListUC(collectionRepo)
	collectionCreateUC := collection.NewCreateUC(collectionRepo)
	collectionAddCardUC := collection.NewAddCardUC(collectionRepo, cardsRepository)
	collectionExportUC := collection.NewExportUC(collectionRepo)

	collectionHandler := collectionHttp.NewHandler(
		collectionListUC,
		collectionCreateUC,
		collectionAddCardUC,
		collectionExportUC,
	)

	aiAgent := deck.NewAIAgent()
	deckRepo := deck.NewRepository(DB)
	deckBuildUC := deck.CreateBuildDeckUC(collectionRepo, aiAgent)
	deckSaveUC := deck.NewSaveUC(deckRepo, cardsRepository)
	deckListUC := deck.NewListUC(deckRepo, imageURIBuilder)
	deckHandler := deckHttp.NewHandler(deckBuildUC, deckSaveUC, deckListUC)

	return &Container{
		CardHandler:       cardHandler,
		CollectionHandler: collectionHandler,
		DeckHandler:       deckHandler,
	}, nil
}
