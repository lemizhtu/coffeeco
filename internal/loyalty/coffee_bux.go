package loyalty

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	coffeeco "github.com/lemizhtu/coffeego/internal"
	"github.com/lemizhtu/coffeego/internal/store"
)

type CoffeeBux struct {
	ID                                    uuid.UUID
	store                                 store.Store
	coffeeLover                           coffeeco.CoffeeLover
	FreeDrinksAvailable                   int
	RemainingDrinkPurchasesUntilFreeDrink int
}

func (cb *CoffeeBux) AddStamp() {
	if cb.RemainingDrinkPurchasesUntilFreeDrink == 1 {
		cb.RemainingDrinkPurchasesUntilFreeDrink = 10
		cb.FreeDrinksAvailable++
	} else {
		cb.RemainingDrinkPurchasesUntilFreeDrink--
	}
}

func (cb *CoffeeBux) Pay(_ context.Context, productsToPurchase []coffeeco.Product) error {
	lp := len(productsToPurchase)

	if lp == 0 {
		return errors.New("nothing to buy")
	}

	if cb.FreeDrinksAvailable < lp {
		return fmt.Errorf("not enough coffeeBux to cover entire purchase. Have %d, need %d", lp, cb.FreeDrinksAvailable)
	}

	cb.FreeDrinksAvailable -= lp

	return nil
}
