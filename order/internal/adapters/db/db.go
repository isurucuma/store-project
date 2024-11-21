package db

import (
	"context"
	"fmt"
	"github.com/isurucuma/store-project/order/internal/application/domain"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	CustomerId int64
	Status     string
	OrderItems []OrderItem
}

type OrderItem struct {
	gorm.Model
	ProductCode string
	UnitPrice   float32
	Quantity    int32
	OrderID     uint
}

type Adapter struct {
	db *gorm.DB
}

func (a Adapter) Get(ctx context.Context, id int64) (domain.Order, error) {
	var orderEntity Order
	res := a.db.WithContext(ctx).Preload("OrderItems").First(&orderEntity, id)
	if res.Error != nil {
		return domain.Order{}, res.Error
	}
	var orderItems []domain.OrderItem
	for _, orderItem := range orderEntity.OrderItems {
		orderItems = append(orderItems, domain.OrderItem{
			ProductCode: orderItem.ProductCode,
			UnitPrice:   orderItem.UnitPrice,
			Quantity:    orderItem.Quantity,
		})
	}

	order := domain.Order{
		ID:         int64(orderEntity.ID),
		CustomerId: orderEntity.CustomerId,
		Status:     domain.OrderStatus(orderEntity.Status),
		OrderItems: orderItems,
		CreatedAt:  orderEntity.CreatedAt.UnixNano(),
	}

	return order, nil
}

func (a Adapter) Save(ctx context.Context, order *domain.Order) error {
	var orderItems []OrderItem
	for _, orderItem := range order.OrderItems {
		orderItems = append(orderItems, OrderItem{
			ProductCode: orderItem.ProductCode,
			UnitPrice:   orderItem.UnitPrice,
			Quantity:    orderItem.Quantity,
		})
	}
	orderModel := Order{
		CustomerId: order.CustomerId,
		Status:     string(order.Status),
		OrderItems: orderItems,
	}
	res := a.db.WithContext(ctx).Create(&orderModel)
	if res.Error == nil {
		order.ID = int64(orderModel.ID)
	}
	return res.Error
}

func NewAdapter(dataSourceUrl string) (*Adapter, error) {
	db, err := gorm.Open(mysql.Open(dataSourceUrl), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("db connection error: %v", err)
	}
	err = db.AutoMigrate(&Order{}, &OrderItem{})
	if err != nil {
		return nil, fmt.Errorf("db migration error: %v", err)
	}
	return &Adapter{db: db}, nil
}
