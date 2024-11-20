package ports

import (
	"context"
	"github.com/isurucuma/store-project/order/internal/application/domain"
)

// APIPort is the interface that has the contract for the functions that our application will be provided to the outside
//
// # These functions are implemented by the Application
//
// If we need to change the implementation of the application, we can easily do it without changing the API. So that the
// system will be more maintainable and testable.
type APIPort interface {
	PlaceOrder(ctx context.Context, order domain.Order) (domain.Order, error)
	GetOrder(ctx context.Context, orderId int64) (domain.Order, error)
}
