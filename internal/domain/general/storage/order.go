package storage

import (
	"errors"

	"github.com/is-web-y26/m3302-milovatskiy/internal/domain/general"
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	UserID uint
	Phones []general.Phone `gorm:"many2many:order_phones;"`
	Status OrderStatus
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
