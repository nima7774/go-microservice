package order

import (
	"context"
	"time"

	"github.com/segmentio/ksuid"
)

type service interface {
	PostOrder(ctx context.Context, accountID string, products []OrderedProduct) (*Order, error)
	GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error)
}

type Order struct {
	ID         string           `json:"id"`
	CreatedAt  time.Time        `json:"created_at"`
	TotalPrice float64          `json:"total_price"`
	AccountID  string           `json:"account_id"`
	Products   []OrderedProduct `json:"products"`
}

type OrderedProduct struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
}

type OrderService struct {
	repository Repository
}

func NewOrderService(r Repository) *OrderService {
	return &OrderService{repository: r}
}

func (s *OrderService) PostOrder(ctx context.Context, accountID string, products []OrderedProduct) (*Order, error) {
	o := &Order{
		ID:        ksuid.New().String(),
		CreatedAt: time.Now().UTC(),

		AccountID: accountID,
		Products:  products,
	}
	o.TotalPrice = 0.0
	for _, p := range products {
		o.TotalPrice += p.Price * float64(p.Quantity)
	}
	err := s.repository.PutOrder(ctx, *o)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (s *OrderService) GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error) {
	return s.repository.GetOrdersForAccount(ctx, accountID)
}
