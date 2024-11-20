package api

import (
	"context"
	"github.com/isurucuma/store-project/order/internal/application/domain"
	"github.com/isurucuma/store-project/order/internal/ports"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Application struct {
	db      ports.DBPort
	payment ports.PaymentPort
}

func NewApplication(db ports.DBPort, payment ports.PaymentPort) *Application {
	return &Application{
		db:      db,
		payment: payment,
	}
}

func (a Application) PlaceOrder(ctx context.Context, order domain.Order) (domain.Order, error) {
	err := a.db.Save(ctx, order)
	if err != nil {
		return domain.Order{}, err
	}
	err = a.payment.Charge(ctx, order)
	if err != nil {
		errorStatus, _ := status.FromError(err)
		fieldError := &errdetails.BadRequest_FieldViolation{
			Field:       "payment",
			Description: errorStatus.Message(),
		}
		badReq := &errdetails.BadRequest{}
		badReq.FieldViolations = append(badReq.FieldViolations, fieldError)
		orderStatus := status.New(codes.InvalidArgument, "order creation failed")
		statusWithDetails, _ := orderStatus.WithDetails(badReq)
		return domain.Order{}, statusWithDetails.Err()
	}
	return order, nil
}

func (a Application) GetOrder(ctx context.Context, orderId int64) (domain.Order, error) {
	return a.db.Get(ctx, orderId)
}
