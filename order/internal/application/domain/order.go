package domain

import "time"

type OrderStatus string

const (
	PENDING    OrderStatus = "PENDING"
	PROCESSING OrderStatus = "PROCESSING"
	CANCELLED  OrderStatus = "CANCELLED"
	COMPLETED  OrderStatus = "COMPLETED"
	SHIPPED    OrderStatus = "SHIPPED"
)

type Order struct {
	ID         int64       `json:"id"`
	CustomerId int64       `json:"customer_id"`
	OrderItems []OrderItem `json:"order_items"`
	Status     OrderStatus `json:"status"`
	CreatedAt  int64       `json:"created_at"`
}

type OrderItem struct {
	ProductCode int64   `json:"product_code"`
	UnitPrice   float32 `json:"unit_price"`
	Quantity    int32   `json:"quantity"`
}

func NewOrder(customerId int64, orderItems []OrderItem) Order {
	return Order{
		CreatedAt:  time.Now().Unix(),
		OrderItems: orderItems,
		Status:     PENDING,
		CustomerId: customerId,
	}
}

func (o *Order) TotalPrice() (total float32) {
	for _, item := range o.OrderItems {
		total += item.UnitPrice * float32(item.Quantity)
	}
	return
}
