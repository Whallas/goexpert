package service

import (
	"context"

	pb "github.com/Whallas/goexpert/Desafios/clean-architecture/internal/infra/grpc/pb"
	"github.com/Whallas/goexpert/Desafios/clean-architecture/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer
	CreateOrderUseCase usecase.CreateOrderUseCase
	ListOrdersUseCase  usecase.ListOrdersUseCase
}

func NewOrderService(create usecase.CreateOrderUseCase, list usecase.ListOrdersUseCase) *OrderService {
	return &OrderService{
		CreateOrderUseCase: create,
		ListOrdersUseCase:  list,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	input := usecase.OrderInputDTO{Price: req.Price, Tax: req.Tax}
	output, err := s.CreateOrderUseCase.Execute(input)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.CreateOrderResponse{
		Id:         output.ID,
		Price:      output.Price,
		Tax:        output.Tax,
		FinalPrice: output.FinalPrice,
	}, nil
}

func (s *OrderService) ListOrders(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	outputs, err := s.ListOrdersUseCase.Execute()
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	orders := make([]*pb.CreateOrderResponse, 0, len(outputs))
	for _, o := range outputs {
		orders = append(orders, &pb.CreateOrderResponse{
			Id:         o.ID,
			Price:      o.Price,
			Tax:        o.Tax,
			FinalPrice: o.FinalPrice,
		})
	}
	return &pb.ListOrdersResponse{Orders: orders}, nil
}
