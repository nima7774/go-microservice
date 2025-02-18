package catalog

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/nima7774/go-microservice/catalog/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	service Service
	pb.UnimplementedCatalogServiceServer
}

func NewGRPCServer(service Service) *grpcServer {
	return &grpcServer{service: service}
}
func ListenGRPC(s Service, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return err
	}
	serv := grpc.NewServer()
	pb.RegisterCatalogServiceServer(serv, &grpcServer{
		UnimplementedCatalogServiceServer: pb.UnimplementedCatalogServiceServer{},
		service:                           s,
	})
	reflection.Register(serv)
	return serv.Serve(lis)
}

func (s *grpcServer) PostProduct(ctx context.Context, r *pb.PostProductRequest) (*pb.PostProductResponse, error) {
	p, err := s.service.PostProduct(ctx, r.Name, r.Description, r.Price)
	if err != nil {
		log.Printf("error posting product: %v", err)
		return nil, err
	}
	return &pb.PostProductResponse{Product: &pb.Product{
		Id:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
	}}, nil
}

func (s *grpcServer) GetProduct(ctx context.Context, r *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	p, err := s.service.GetProduct(ctx, r.Id)
	if err != nil {
		log.Printf("error getting product: %v", err)
		return nil, err
	}

	return &pb.GetProductResponse{Product: &pb.Product{
		Id:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
	}}, nil
}

func (s *grpcServer) GetProducts(ctx context.Context, r *pb.GetProductsRequest) (*pb.GetProductsResponse, error) {
	var res []Product
	var err error
	if r.Query != "" {
		res, err = s.service.SearchProducts(ctx, r.Query, r.Skip, r.Take)
	} else if len(r.Ids) > 0 {
		res, err = s.service.GetProductsByIDs(ctx, r.Ids)
	} else {
		res, err = s.service.GetProducts(ctx, r.Skip, r.Take)
	}
	if err != nil {
		log.Printf("error getting products: %v", err)
		return nil, err
	}
	products := []*pb.Product{}
	for _, p := range res {
		products = append(products, &pb.Product{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		})
	}
	return &pb.GetProductsResponse{Products: products}, nil
}
