package payment

import (
	"context"
	"errors"
	"fmt"

	"github.com/Rhymond/go-money"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/charge"
	"github.com/stripe/stripe-go/v76/client"
)

type StripeService struct {
	stripeClient *client.API
}

func NewStripeService(apiKey string) (*StripeService, error) {
	if apiKey == "" {
		return nil, errors.New("API key cannot be empty")
	}

	var (
		sc = &client.API{}
	)

	sc.Init(apiKey, nil)

	return &StripeService{
		stripeClient: sc,
	}, nil
}

func (s StripeService) ChargeCard(_ context.Context, amount money.Money, cardToken string) error {
	params := &stripe.ChargeParams{
		Amount:   stripe.Int64(amount.Amount()),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		Source:   &stripe.PaymentSourceSourceParams{Token: stripe.String(cardToken)},
	}
	_, err := charge.New(params)
	if err != nil {
		return fmt.Errorf("failed to create a charge: %w", err)
	}

	return nil
}