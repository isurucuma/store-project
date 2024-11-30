package ports

import (
	"context"
	"github.com/isurucuma/store-project/order/internal/application/domain"
)

// PaymentPort is the interface that has the contract for the functions that will be implemented by the PaymentAdapter.
//
// Our application depends on this interface and not on the actual implementation.
//
// Therefore, we can easily switch the implementation of the underline payment service without changing the application.
type PaymentPort interface {
	Charge(ctx context.Context, order *domain.Order) error
}
