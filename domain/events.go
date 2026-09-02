package domain

import "time"

const OrderCreatedEvent = "warehouse.order.created"

type OrderCreated struct {
	EventID   string     `json:"event_id"`
	EventType string     `json:"event_type"`
	Occurred  time.Time  `json:"occurred_at"`
	Order     OrderEvent `json:"order"`
}

type OrderEvent struct {
	ID         uint             `json:"id"`
	CustomerID uint             `json:"customer_id"`
	Status     string           `json:"status"`
	Total      uint             `json:"total"`
	Items      []OrderItemEvent `json:"items"`
}

type OrderItemEvent struct {
	ProductID uint `json:"product_id"`
	Quantity  uint `json:"quantity"`
	UnitPrice uint `json:"unit_price"`
}
