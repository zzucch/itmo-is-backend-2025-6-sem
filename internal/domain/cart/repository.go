package cart

import (
	"github.com/is-web-y26/m3302-milovatskiy/internal/domain/general"
	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateOrder(order *Order) error
	GetPhonesByIDs(phoneIDs []uint) ([]general.Phone, error)
}

type GormOrderRepository struct {
	db *gorm.DB
}

func NewGormOrderRepository(db *gorm.DB) *GormOrderRepository {
	return &GormOrderRepository{db: db}
}

func (r *GormOrderRepository) CreateOrder(order *Order) error {
	return r.db.Create(order).Error
}
