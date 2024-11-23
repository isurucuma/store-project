package grpc

import (
	"context"
	"fmt"
	"github.com/isurucuma/store-project/order/config"
	"github.com/isurucuma/store-project/order/internal/application/domain"
	"github.com/isurucuma/store-project/order/internal/ports"
	"github.com/isurucuma/store-proto/golang/order"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

// Adapter is the gRPC server adapter that will bind the gRPC methods and the APIPort implementation
type Adapter struct {
	server *grpc.Server
	api    ports.APIPort
	port   int
	order.UnimplementedOrderServer
}

func NewAdapter(api ports.APIPort, port int) *Adapter {
	return &Adapter{
		port: port,
		api:  api,
	}
}

func (a Adapter) Run() {
	var err error
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		log.Fatalf("failed to listen on port %d, error: %v", a.port, err)
	}

	grpcServer := grpc.NewServer()
	a.server = grpcServer

	order.RegisterOrderServer(grpcServer, a)
	if config.GetEnv("ENV", "dev") == "dev" {
		// Register reflection service on gRPC server.
		// This is needed for grpcurl to work.
		reflection.Register(grpcServer)
	}

	log.Printf("starting order service on port %d ...", a.port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve grpc on port %d", a.port)
	}
}

func (a Adapter) Stop() {
	a.server.Stop()
}

// these server functions will be invoked by the gRPC server and we should call the APIPort implementation

func (a Adapter) Create(ctx context.Context, request *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {
	var orderItems []domain.OrderItem
	for _, orderItem := range request.OrderItems {
		orderItems = append(orderItems, domain.OrderItem{
			ProductCode: orderItem.ProductCode,
			UnitPrice:   orderItem.UnitPrice,
			Quantity:    orderItem.Quantity,
		})
	}

	newOrder := domain.NewOrder(request.UserId, orderItems)
	result, err := a.api.PlaceOrder(ctx, newOrder)
	if err != nil {
		return nil, err
	}
	return &order.CreateOrderResponse{OrderId: result.ID}, nil
}

func (a Adapter) Get(ctx context.Context, request *order.GetOrderRequest) (*order.GetOrderResponse, error) {
	result, err := a.api.GetOrder(ctx, request.OrderId)
	if err != nil {
		return nil, err
	}
	var orderItems []*order.OrderItem
	for _, orderItem := range result.OrderItems {
		orderItems = append(orderItems, &order.OrderItem{
			ProductCode: orderItem.ProductCode,
			UnitPrice:   orderItem.UnitPrice,
			Quantity:    orderItem.Quantity,
		})
	}
	return &order.GetOrderResponse{UserId: result.CustomerId, OrderItems: orderItems}, nil
}
