package purchase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	coffeeco "github.com/lemizhtu/coffeego/internal"
	"github.com/lemizhtu/coffeego/internal/payment"
	"github.com/lemizhtu/coffeego/internal/store"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository interface {
	Store(ctx context.Context, purchase Purchase) error
}

type MongoRepository struct {
	purchases *mongo.Collection
}

func NewMongoRepo(ctx context.Context, uri string) (*MongoRepository, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to create a mongo client: %w", err)
	}

	var (
		purchases = client.Database("coffeeco").Collection("purchases")
	)

	return &MongoRepository{
		purchases: purchases,
	}, nil
}

func (mr *MongoRepository) Store(ctx context.Context, purchase Purchase) error {
	var (
		p = toMongoPurchase(purchase)
	)

	_, err := mr.purchases.InsertOne(ctx, p)
	if err != nil {
		return fmt.Errorf("failed to persist purchase: %w", err)
	}

	return nil
}

type mongoPurchase struct {
	ID                 uuid.UUID          `bson:"id"`
	Store              store.Store        `bson:"store"`
	ProductsToPurchase []coffeeco.Product `bson:"products_to_purchase"`
	Total              int64              `bson:"total"`
	PaymentMeans       payment.Means      `bson:"payment_means"`
	TimeOfPurchase     time.Time          `bson:"time_of_purchase"`
	CardToken          *string            `bson:"card_token"`
}

func toMongoPurchase(p Purchase) mongoPurchase {
	return mongoPurchase{
		ID:                 p.ID,
		Store:              p.Store,
		ProductsToPurchase: p.ProductsToPurchase,
		Total:              p.total.Amount(),
		PaymentMeans:       p.PaymentMeans,
		TimeOfPurchase:     p.timeOfPurchase,
		CardToken:          p.CardToken,
	}
}
