package cart

import (
	"errors"
)

type Service struct {
	repository Repository
}

func NewService(repo Repository) *Service {
	return &Service{repository: repo}
}

func (s *Service) PlaceOrder(userID uint, phoneIDs []uint) (*Order, error) {
	if len(phoneIDs) == 0 {
		return nil, errors.New("no phones selected")
	}

	phones, err := s.repository.GetPhonesByIDs(phoneIDs)
	if err != nil {
		return nil, err
	}

	order := &Order{
		UserID: userID,
		Phones: phones,
		Status: Pending,
	}

	if err := s.repository.CreateOrder(order); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *Service) GetOrdersByUserID(userID uint) (Order, error) {
	return s.repository.FindOrderByUserID(userID)
}

func (s *Service) DeleteOrder(id uint) error {
	return s.repository.DeleteOrderByID(id)
}

func (s *Service) UpdateOrder(order *Order) (Order, error) {
	if err := s.repository.UpdateOrder(order); err != nil {
		return Order{}, err
	}
	return *order, nil
}

func (s *Service) GetOrderByID(id uint) (Order, error) {
	return s.repository.FindOrderByID(id)
}

func (s *Service) GetAllOrders() ([]Order, error) {
	return s.repository.FindAllOrders()
}
