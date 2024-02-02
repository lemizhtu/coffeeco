package purchase

import (
	"context"
	"fmt"
	"time"

	"github.com/Rhymond/go-money"
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
	id                 uuid.UUID
	store              store.Store
	productsToPurchase []coffeeco.Product
	total              money.Money
	paymentMeans       payment.Means
	timeOfPurchase     time.Time
	cardToken          *string
}

func toMongoPurchase(p Purchase) mongoPurchase {
	return mongoPurchase{
		id:                 p.ID,
		store:              p.Store,
		productsToPurchase: p.ProductsToPurchase,
		total:              p.total,
		paymentMeans:       p.PaymentMeans,
		timeOfPurchase:     p.timeOfPurchase,
		cardToken:          p.CardToken,
	}
}
