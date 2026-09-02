package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"goproject/domain"
	"goproject/domain/dto"
	"goproject/repository"
)

type OrderUsecase interface {
	Create(ctx context.Context, request dto.CreateOrderRequest) (*domain.Order, error)
}

type orderUsecase struct {
	orders    repository.OrderRepository
	products  repository.ProductRepository
	publisher repository.EventPublisher
	topic     string
	timeout   time.Duration
}

func NewOrderUsecase(orders repository.OrderRepository, products repository.ProductRepository, publisher repository.EventPublisher, topic string, timeout time.Duration) OrderUsecase {
	return &orderUsecase{orders: orders, products: products, publisher: publisher, topic: topic, timeout: timeout}
}

func (service *orderUsecase) Create(parent context.Context, request dto.CreateOrderRequest) (*domain.Order, error) {
	ctx, cancel := context.WithTimeout(parent, service.timeout)
	defer cancel()
	order := &domain.Order{CustomerID: request.CustomerID, Status: "created"}
	for _, item := range request.Items {
		product, err := service.products.GetByID(ctx, item.ProductID)
		if err != nil {
			return nil, err
		}
		order.Items = append(order.Items, domain.OrderItem{ProductID: product.ID, Quantity: item.Quantity, UnitPrice: product.Price})
		order.Total += product.Price * item.Quantity
	}
	if err := service.orders.Create(ctx, order); err != nil {
		return nil, err
	}
	if service.publisher == nil {
		return nil, errors.New("order event publisher is unavailable; order was not accepted")
	}
	event, err := newOrderCreatedEvent(order)
	if err != nil {
		return nil, err
	}
	payload, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}
	if err := service.publisher.Publish(ctx, service.topic, repository.Event{Key: []byte(stringID(order.ID)), Value: payload}); err != nil {
		return nil, err
	}
	return order, nil
}

func newOrderCreatedEvent(order *domain.Order) (domain.OrderCreated, error) {
	eventID := make([]byte, 16)
	if _, err := rand.Read(eventID); err != nil {
		return domain.OrderCreated{}, err
	}
	items := make([]domain.OrderItemEvent, 0, len(order.Items))
	for _, item := range order.Items {
		items = append(items, domain.OrderItemEvent{ProductID: item.ProductID, Quantity: item.Quantity, UnitPrice: item.UnitPrice})
	}
	return domain.OrderCreated{EventID: hex.EncodeToString(eventID), EventType: domain.OrderCreatedEvent, Occurred: time.Now().UTC(), Order: domain.OrderEvent{ID: order.ID, CustomerID: order.CustomerID, Status: order.Status, Total: order.Total, Items: items}}, nil
}

func stringID(id uint) string {
	return strconv.FormatUint(uint64(id), 10)
}
