package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrNoDiscount = errors.New("no discount for store")
)

type Repository interface {
	GetStoreDiscount(ctx context.Context, storeID uuid.UUID) (float32, error)
}

type MongoRepository struct {
	storeDiscounts *mongo.Collection
}

func NewMongoRepo(ctx context.Context, uri string) (*MongoRepository, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to create a mongo client: %w", err)
	}

	var (
		storeDiscounts = client.Database("coffeeco").Collection("store_discounts")
	)

	return &MongoRepository{
		storeDiscounts: storeDiscounts,
	}, nil
}

func (mr *MongoRepository) GetStoreDiscount(ctx context.Context, storeID uuid.UUID) (float32, error) {
	var (
		discount float32
	)

	if err := mr.storeDiscounts.FindOne(ctx, bson.D{{"store_id", storeID.String()}}).Decode(&discount); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return 0.0, ErrNoDiscount
		}

		return 0.0, fmt.Errorf("failed to find discount for store: %w", err)
	}

	return discount, nil
}
