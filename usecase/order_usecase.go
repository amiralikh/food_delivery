package usecase

import (
	"foodDelivery/domain"
	"foodDelivery/repository"
)

type OrderUseCase interface {
	SubmitOrder(order *domain.Order) error
}

type orderUseCase struct {
	orderRepository repository.OrderRepository
}

func NewOrderUseCase(orderRepository repository.OrderRepository) OrderUseCase {
	return &orderUseCase{
		orderRepository: orderRepository,
	}
}

func (ou *orderUseCase) SubmitOrder(order *domain.Order) error {
	if len(*order.Items) == 0 {
		return repository.ErrItemsNotFound
	}
	err := ou.orderRepository.SubmitOrder(order)
	if err != nil {
		return err
	}
	return nil
}
