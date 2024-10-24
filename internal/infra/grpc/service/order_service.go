package service

import (
	"context"
	"log"

	"github.com/devfullcycle/20-CleanArch/internal/entity"
	"github.com/devfullcycle/20-CleanArch/internal/infra/grpc/pb"
	"github.com/devfullcycle/20-CleanArch/internal/usecase"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer
	CreateOrderUseCase usecase.CreateOrderUseCase
	OrderRepository    entity.OrderRepositoryInterface
}

func NewOrderService(createOrderUseCase usecase.CreateOrderUseCase) *OrderService {
	return &OrderService{
		CreateOrderUseCase: createOrderUseCase,
		OrderRepository:    createOrderUseCase.OrderRepository,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, in *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	log.Printf("Executing CreateOrderUseCase with repository: %+v", s.CreateOrderUseCase.OrderRepository)

	dto := usecase.OrderInputDTO{
		ID:    in.Id,
		Price: float64(in.Price),
		Tax:   float64(in.Tax),
	}
	output, err := s.CreateOrderUseCase.Execute(dto)
	if err != nil {
		return nil, err
	}
	return &pb.CreateOrderResponse{
		Id:         output.ID,
		Price:      float32(output.Price),
		Tax:        float32(output.Tax),
		FinalPrice: float32(output.FinalPrice),
	}, nil
}

func (s *OrderService) ListOrders(ctx context.Context, in *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	// Use the use case to fetch orders
	listOrdersUseCase := usecase.NewListOrdersUseCase(s.OrderRepository)
	orders, err := listOrdersUseCase.Execute()
	if err != nil {
		return nil, err
	}

	// Convert Go orders to pb.Order for gRPC response
	var pbOrders []*pb.Order
	for _, order := range orders {
		pbOrders = append(pbOrders, &pb.Order{
			Id:         order.ID,
			Price:      float32(order.Price),
			Tax:        float32(order.Tax),
			FinalPrice: float32(order.FinalPrice),
		})
	}

	// Return the ListOrdersResponse with the list of orders
	return &pb.ListOrdersResponse{Orders: pbOrders}, nil
}
