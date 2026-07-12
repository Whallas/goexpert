package usecase

import "github.com/Whallas/goexpert/Desafios/clean-architecture/internal/domain/entity"

type ListOrdersUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
}

func NewListOrdersUseCase(repo entity.OrderRepositoryInterface) *ListOrdersUseCase {
	return &ListOrdersUseCase{OrderRepository: repo}
}

func (l *ListOrdersUseCase) Execute() ([]OrderOutputDTO, error) {
	orders, err := l.OrderRepository.GetAll()
	if err != nil {
		return nil, err
	}
	output := make([]OrderOutputDTO, 0, len(orders))
	for _, o := range orders {
		output = append(output, OrderOutputDTO{
			ID:         o.ID,
			Price:      o.Price,
			Tax:        o.Tax,
			FinalPrice: o.FinalPrice,
		})
	}
	return output, nil
}
