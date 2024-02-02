package main

import (
	"context"
	"log"

	"github.com/Rhymond/go-money"
	"github.com/google/uuid"
	coffeeco "github.com/lemizhtu/coffeego/internal"
	"github.com/lemizhtu/coffeego/internal/payment"
	"github.com/lemizhtu/coffeego/internal/purchase"
	"github.com/lemizhtu/coffeego/internal/store"
)

func main() {
	var (
		ctx          = context.Background()
		stripeAPIKey = "sk_test_4eC39HqLyjWDarjtT1zdp7dc"
		cardToken    = "tok_visa"
		mongoURI     = "mongodb://root:root@localhost:27017"
	)

	stripeService, err := payment.NewStripeService(stripeAPIKey)
	if err != nil {
		log.Fatal(err)
	}

	purchaseMongoRepo, err := purchase.NewMongoRepo(ctx, mongoURI)
	if err != nil {
		log.Fatal(err)
	}

	storeMongoRepo, err := store.NewMongoRepo(ctx, mongoURI)
	if err != nil {
		log.Fatal(err)
	}

	var (
		storeService    = store.NewService(storeMongoRepo)
		purchaseService = purchase.NewService(stripeService, purchaseMongoRepo, storeService)
		someStoreID     = uuid.New()
	)

	pur := purchase.Purchase{
		Store: store.Store{
			ID: someStoreID,
		},
		ProductsToPurchase: []coffeeco.Product{
			{
				ItemName:  "Item1",
				BasePrice: *money.New(3000, money.USD),
			},
		},
		PaymentMeans: payment.MeansCard,
		CardToken:    &cardToken,
	}
	if err := purchaseService.CompletePurchase(ctx, someStoreID, &pur, nil); err != nil {
		log.Fatal(err)
	}

	log.Println("purchase was successful")
}
