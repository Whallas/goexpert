package usecase

import (
	"github.com/Whallas/goexpert/Desafios/clean-architecture/internal/domain/entity"
	"github.com/google/uuid"
)

type OrderInputDTO struct {
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
}

func NewCreateOrderUseCase(repo entity.OrderRepositoryInterface) *CreateOrderUseCase {
	return &CreateOrderUseCase{OrderRepository: repo}
}

func (c *CreateOrderUseCase) Execute(input OrderInputDTO) (OrderOutputDTO, error) {
	order, err := entity.NewOrder(uuid.New().String(), input.Price, input.Tax)
	if err != nil {
		return OrderOutputDTO{}, err
	}
	order.CalculateFinalPrice()
	if err := c.OrderRepository.Save(order); err != nil {
		return OrderOutputDTO{}, err
	}
	return OrderOutputDTO{
		ID:         order.ID,
		Price:      order.Price,
		Tax:        order.Tax,
		FinalPrice: order.FinalPrice,
	}, nil
}
