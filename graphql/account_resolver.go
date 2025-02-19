package main

import (
	"context"
	"log"
	"time"
)

type accountResolver struct {
	server *Server
}

func (r *accountResolver) Orders(ctx context.Context, obj *Account) ([]*Order, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	orderList, err := r.server.orderClient.GetOrders(ctx, obj.ID)
	if err != nil {
		log.Printf("Error getting orders: %v", err)
		return nil, err
	}

	var orders []*Order
	for _, o := range orderList {
		var products []*OrderedProduct
		for _, p := range o.Products {
			products = append(products, &OrderedProduct{
				ID:          p.ID,
				Name:        p.Name,
				Price:       p.Price,
				Description: p.Description,
				Quantity:    int(p.Quantity),
			})
		}
		orders = append(orders, &Order{
			ID:         o.ID,
			Products:   products,
			CreatedAt:  o.CreatedAt,
			TotalPrice: o.TotalPrice,
		})
	}

	return nil, nil
}
