package usecase

import (
	"fmt"
	"log"

	"github.com/devfullcycle/20-CleanArch/internal/entity"
	"github.com/devfullcycle/20-CleanArch/pkg/events"
)

type OrderInputDTO struct {
	ID    string  `json:"id"`
	Price float64 `json:"price"`
	Tax   float64 `json:"tax"`
}

type OrderOutputDTO struct {
	ID         string  `json:"id"`
	Price      float64 `json:"price"`
	Tax        float64 `json:"tax"`
	FinalPrice float64 `json:"final_price"`
}

type CreateOrderUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
	OrderCreated    events.EventInterface
	EventDispatcher events.EventDispatcherInterface
}

func NewCreateOrderUseCase(
	OrderRepository entity.OrderRepositoryInterface,
	OrderCreated events.EventInterface,
	EventDispatcher events.EventDispatcherInterface,
) *CreateOrderUseCase {
	fmt.Printf("Creating CreateOrderUseCase with OrderRepository: %+v, OrderCreated: %+v, EventDispatcher: %+v\n",
		OrderRepository, OrderCreated, EventDispatcher)

	if OrderRepository == nil {
		panic("OrderRepository is nil")
	}
	return &CreateOrderUseCase{
		OrderRepository: OrderRepository,
		OrderCreated:    OrderCreated,
		EventDispatcher: EventDispatcher,
	}
}

func (c *CreateOrderUseCase) Execute(input OrderInputDTO) (OrderOutputDTO, error) {
	// Print the repository to confirm if it is non-nil at execution time
	fmt.Printf("Executing CreateOrderUseCase with repository: %+v\n", c.OrderRepository)
	log.Printf("Inside Execute with OrderRepository: %+v", c.OrderRepository)

	if c.OrderRepository == nil {
		return OrderOutputDTO{}, fmt.Errorf("OrderRepository is not initialized")
	}

	// Debug statements to check for nil dependencies
	if c.OrderRepository == nil {
		return OrderOutputDTO{}, fmt.Errorf("OrderRepository is not initialized")
	}
	if c.OrderCreated == nil {
		return OrderOutputDTO{}, fmt.Errorf("OrderCreated event is not initialized")
	}
	if c.EventDispatcher == nil {
		return OrderOutputDTO{}, fmt.Errorf("EventDispatcher is not initialized")
	}

	order := entity.Order{
		ID:    input.ID,
		Price: input.Price,
		Tax:   input.Tax,
	}
	order.CalculateFinalPrice()
	if err := c.OrderRepository.Save(&order); err != nil {
		return OrderOutputDTO{}, err
	}

	dto := OrderOutputDTO{
		ID:         order.ID,
		Price:      order.Price,
		Tax:        order.Tax,
		FinalPrice: order.Price + order.Tax,
	}

	c.OrderCreated.SetPayload(dto)
	c.EventDispatcher.Dispatch(c.OrderCreated)

	return dto, nil
}
