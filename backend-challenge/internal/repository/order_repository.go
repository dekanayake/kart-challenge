package repository

import (
	"time"

	"github.com/dekanayake/kart-challenge/backend-challenge/internal/config"
	"github.com/google/uuid"
)

type OrderRepository interface {
	CreateOrder(items []OrderItem, couponCode string) (*Order, error)
}

type InMemoryOrderRepository struct{}

func newInMemoryOrderRepository() *InMemoryOrderRepository {
	return &InMemoryOrderRepository{}
}

func (r *InMemoryOrderRepository) CreateOrder(items []OrderItem, couponCode string) (*Order, error) {
	order := Order{
		ID:        uuid.NewString(),
		Items:     items,
		CreatedAt: time.Now(),
	}

	config.Logger.Info().
		Str("order_id", order.ID).
		Int("items_count", len(items)).
		Interface("items", items).
		Time("created_at", order.CreatedAt).
		Msg("New order created")

	return &order, nil
}
