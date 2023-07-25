package usecase

import (
	"foodDelivery/domain"
	"foodDelivery/repository"
)

type OrderUseCase interface {
	SubmitOrder(order *domain.Order) error
	GetUserOrders(userId int64) (*[]domain.Order, error)
	GetOrderWithItems(orderID int64) (*domain.Order, error)
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

func (ou *orderUseCase) GetUserOrders(userId int64) (*[]domain.Order, error) {
	orders, err := ou.orderRepository.GetUserOrders(userId)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (ou *orderUseCase) GetOrderWithItems(orderID int64) (*domain.Order, error) {
	order, err := ou.orderRepository.GetOrderWithItems(orderID)
	if err != nil {
		return nil, err
	}
	return order, nil
}
