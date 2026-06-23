// Package user : repository_mongo.go Describe mongo impletation of user data store
package user

import (
	"context"
	"errors"
	"time"

	"github.com/meplos/locana-deck-builder/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoRepository struct {
	col *mongo.Collection
}

func NewRepository(DB *mongo.Database) Repository {
	return &MongoRepository{
		col: DB.Collection("users"),
	}
}

func (r *MongoRepository) Create(ctx context.Context, u domain.User) error {
	_, err := r.col.InsertOne(ctx, u)
	return err
}

func (r *MongoRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	err := r.col.FindOne(ctx, bson.D{{Key: "email", Value: email}}).Err()
	if errors.Is(err, mongo.ErrNoDocuments) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *MongoRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	var u domain.User
	err := r.col.FindOne(ctx, bson.D{{Key: "email", Value: email}}).Decode(&u)
	return u, err
}

func (r *MongoRepository) LogUserAt(ctx context.Context, ID string, iat time.Time) error {
	_, err := r.col.UpdateByID(ctx, ID, bson.D{
		{
			Key: "$set",
			Value: bson.M{
				"last_login": iat,
			},
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *MongoRepository) FindByID(ctx context.Context, ID string) (domain.User, error) {
	var u domain.User
	err := r.col.FindOne(ctx, bson.D{{Key: "_id", Value: ID}}).Decode(&u)
	return u, err
}
