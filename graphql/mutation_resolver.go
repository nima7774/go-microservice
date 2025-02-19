package main

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/nima7774/go-microservice/order"
)

var (
	ErrInvalidParameter = errors.New("invalid input")
)

type mutationResolver struct {
	server *Server
}

func (r *mutationResolver) CreateAccount(ctx context.Context, in AccountInput) (*Account, error) {

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	account, err := r.server.accountClient.PostAccount(ctx, in.Name)
	if err != nil {
		log.Printf("Error creating account: %v", err)
		return nil, err
	}

	return &Account{
		ID:   account.ID,
		Name: account.Name,
	}, nil
}
func (r *mutationResolver) CreateProduct(ctx context.Context, in ProductInput) (*Product, error) {

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	p, err := r.server.catalogClient.PostProduct(ctx, in.Name, in.Description, in.Price)
	if err != nil {
		log.Printf("Error creating product: %v", err)
		return nil, err
	}

	return &Product{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
	}, nil
}
func (r *mutationResolver) CreateOrder(ctx context.Context, in OrderInput) (*Order, error) {

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var products []order.OrderedProduct
	for _, p := range in.Products {
		if p.Quantity <= 0 {
			return nil, ErrInvalidParameter
		}
		products = append(products, order.OrderedProduct{
			ID:       p.ID,
			Quantity: int(p.Quantity),
		})
	}

	o, err := r.server.orderClient.PostOrder(ctx, in.AccountID, products)
	if err != nil {
		log.Printf("Error creating order: %v", err)
		return nil, err
	}

	return &Order{
		ID:         o.ID,
		TotalPrice: o.TotalPrice,
		CreatedAt:  o.CreatedAt,
	}, nil
}

// CreateProduct

// CreateOrder
