package cart

import (
	"github.com/is-web-y26/m3302-milovatskiy/internal/domain/general"
)

type Repository interface {
	CreateOrder(order *Order) error
	GetPhonesByIDs(phoneIDs []uint) ([]general.Phone, error)
	FindOrderByUserID(userID uint) (Order, error)
	FindOrderByID(id uint) (Order, error)
	DeleteOrderByID(id uint) error
	UpdateOrder(order *Order) error
	FindAllOrders() ([]Order, error)
}
