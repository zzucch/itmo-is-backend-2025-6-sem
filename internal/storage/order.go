package storage

import (
	"errors"

	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	UserID   uint
	PhoneIDs []uint64
	Status   OrderStatus
}

type OrderStatus string

const (
	Pending   OrderStatus = "Pending"
	Shipped   OrderStatus = "Shipped"
	Delivered OrderStatus = "Delivered"
	Cancelled OrderStatus = "Cancelled"
)

func ValidateOrderStatus(status OrderStatus) error {
	switch status {
	case Pending, Shipped, Delivered, Cancelled:
		return nil
	default:
		return errors.New("invalid order status")
	}
}
