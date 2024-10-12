package graph

import (
	"github.com/devfullcycle/20-CleanArch/internal/entity"
	"github.com/devfullcycle/20-CleanArch/internal/usecase"
)

// Resolver struct that gqlgen will use for dependency injection
type Resolver struct {
	CreateOrderUseCase usecase.CreateOrderUseCase
	ListOrdersUseCase  usecase.ListOrdersUseCase
	OrderRepository    entity.OrderRepositoryInterface // <-- Add OrderRepository
}

// Implement the Query method to return the QueryResolver
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

// queryResolver will handle the queries (like ListOrders)
type queryResolver struct {
	*Resolver
}
