package order

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	account "github.com/nima7774/go-microservice/account"
	catalog "github.com/nima7774/go-microservice/catalog"
	pb "github.com/nima7774/go-microservice/order/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	pb.UnimplementedOrderServiceServer
	service       Service
	accountClient *account.Client
	catalogClient *catalog.Client
}

func ListenGRPC(s Service, AccountURL string, catalogURL string, port int) error {
	accountClient, err := account.NewClient(AccountURL)
	if err != nil {
		return err
	}

	catalogClient, err := catalog.NewClient(catalogURL)
	if err != nil {
		accountClient.Close()
		return err
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		catalogClient.Close()
		accountClient.Close()
		return err
	}

	serv := grpc.NewServer()
	pb.RegisterOrderServiceServer(serv, &grpcServer{
		UnimplementedOrderServiceServer: pb.UnimplementedOrderServiceServer{},
		service:                         s,
		accountClient:                   accountClient,
		catalogClient:                   catalogClient,
	})

	reflection.Register(serv)

	return serv.Serve(lis)
}

func (s *grpcServer) PostOrder(ctx context.Context, r *pb.PostOrderRequest) (*pb.PostOrderResponse, error) {
	_, err := s.accountClient.GetAccount(ctx, r.AccountId)
	if err != nil {
		log.Println("error getting account", err)
		return nil, errors.New("account not found")
	}

	productIDs := []string{}

	orderedProducts, err := s.catalogClient.GetProducts(ctx, 0, 0, productIDs, "")
	if err != nil {
		log.Println("error getting products", err)
		return nil, errors.New("products not found")
	}
	products := []OrderedProduct{}
	for _, p := range orderedProducts {
		product := OrderedProduct{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Quantity:    0,
		}
		for _, rp := range r.Products {
			if rp.ProductId == p.ID {
				product.Quantity = int(rp.Quantity)
				break
			}
		}
		if product.Quantity != 0 {
			products = append(products, product)
		}
	}

	order, err := s.service.PostOrder(ctx, r.AccountId, products)
	if err != nil {
		log.Println("error posting order", err)
		return nil, errors.New("error posting order")
	}

	orderProto := &pb.Order{
		Id:         order.ID,
		TotalPrice: order.TotalPrice,
		AccountId:  order.AccountID,
		Products:   []*pb.Order_OrderProduct{},
	}
	orderProto.CreatedAt, _ = order.CreatedAt.MarshalBinary()

	for _, p := range order.Products {
		orderProto.Products = append(orderProto.Products, &pb.Order_OrderProduct{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Quantity:    uint32(p.Quantity),
		})
	}

	return &pb.PostOrderResponse{
		Order: orderProto,
	}, nil
}

func (s *grpcServer) GetOrdersForAccount(ctx context.Context, r *pb.GetOrdersForAccountRequest) (*pb.GetOrdersForAccountResponse, error) {
	accountOrders, err := s.service.GetOrdersForAccount(ctx, r.AccountId)
	if err != nil {
		log.Println("error getting orders", err)
		return nil, errors.New("orders not found")
	}
	productIDMap := map[string]bool{}
	for _, o := range accountOrders {
		for _, p := range o.Products {
			productIDMap[p.ID] = true
		}
	}
	productIDs := []string{}

	for id := range productIDMap {
		productIDs = append(productIDs, id)
	}

	products, err := s.catalogClient.GetProducts(ctx, 0, 0, productIDs, "")
	if err != nil {
		log.Println("error getting products", err)
		return nil, errors.New("products not found")
	}

	orders := []*pb.Order{}

	for _, o := range accountOrders {

		op := &pb.Order{
			Id:         o.ID,
			TotalPrice: o.TotalPrice,
			AccountId:  o.AccountID,
			Products:   []*pb.Order_OrderProduct{},
		}
		op.CreatedAt, _ = o.CreatedAt.MarshalBinary()

		for _, product := range o.Products {
			for _, p := range products {
				if p.ID == product.ID {
					product.Name = p.Name
					product.Description = p.Description
					product.Price = p.Price
					product.Quantity = product.Quantity
					break
				}
			}
			op.Products = append(op.Products, &pb.Order_OrderProduct{
				Id:          product.ID,
				Name:        product.Name,
				Description: product.Description,
				Price:       product.Price,
				Quantity:    uint32(product.Quantity),
			})
		}
		orders = append(orders, op)

	}
	return &pb.GetOrdersForAccountResponse{
		Orders: orders,
	}, nil
}
