package ports

import (
	"context"
	"github.com/isurucuma/store-project/order/internal/application/domain"
)

// DBPort is the interface that has the contract for the functions that will be implemented by the DBAdapter.
//
// Our application depends on this interface and not on the actual implementation.
//
// Therefore, we can easily switch the implementation of the underline database without changing the application.
type DBPort interface {
	Get(ctx context.Context, orderId int64) (domain.Order, error)
	Save(ctx context.Context, order domain.Order) error
}
