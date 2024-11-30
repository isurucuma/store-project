package main

import (
	"github.com/isurucuma/store-project/order/config"
	"github.com/isurucuma/store-project/order/internal/adapters/db"
	"github.com/isurucuma/store-project/order/internal/adapters/grpc"
	"github.com/isurucuma/store-project/order/internal/adapters/payment"
	"github.com/isurucuma/store-project/order/internal/application/api"
	"log"
)

func main() {
	dbAdapter, err := db.NewAdapter(config.GetEnv("DATA_SOURCE_URL", "root:admin@tcp(127.0.0.1:3306)/store"))
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	paymentAdapter, err := payment.NewAdapter(config.GetEnv("PAYMENT_SERVICE_URL", ""))
	if err != nil {
		log.Fatalf("failed to connect to payment service: %v", err)
	}

	application := api.NewApplication(dbAdapter, paymentAdapter)
	grpcAdapter := grpc.NewAdapter(application, config.GetApplicationPort())
	grpcAdapter.Run()
}
