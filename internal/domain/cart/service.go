package cart

import (
	"errors"
)

type Service struct {
	repository OrderRepository
}

func NewService(repo OrderRepository) *Service {
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

func (s *Service) GetAllOrders() (any, error) {
	panic("TODO")
}
