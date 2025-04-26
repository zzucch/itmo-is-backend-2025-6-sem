package resolvers

import (
	"strconv"
	"time"

	"github.com/is-web-y26/m3302-milovatskiy/internal/domain/cart"
	"github.com/is-web-y26/m3302-milovatskiy/internal/domain/catalog"
	"github.com/is-web-y26/m3302-milovatskiy/internal/domain/general"
	"github.com/is-web-y26/m3302-milovatskiy/internal/domain/graph/model"
	"github.com/is-web-y26/m3302-milovatskiy/internal/domain/sell"
	"github.com/is-web-y26/m3302-milovatskiy/internal/domain/user"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	PhoneService   *sell.Service
	UserService    *user.Service
	CartService    *cart.Service
	CatalogService *catalog.Service
}

func convertToModelOrder(order *cart.Order) *model.Order {
	var phoneModels []*model.Phone
	for _, phone := range order.Phones {
		phoneModels = append(phoneModels, &model.Phone{
			ID:          strconv.FormatUint(uint64(phone.ID), 10),
			Brand:       phone.Brand,
			Model:       phone.Name,
			Price:       phone.Price,
			Description: phone.Description,
		})
	}

	var userModel *model.User
	if order.UserID != 0 {
		userModel = &model.User{
			ID: strconv.FormatUint(uint64(order.UserID), 10),
		}
	}

	return &model.Order{
		ID:        strconv.FormatUint(uint64(order.ID), 10),
		User:      userModel,
		Phones:    phoneModels,
		Status:    model.OrderStatus(order.Status),
		CreatedAt: order.CreatedAt.Format(time.RFC3339),
		UpdatedAt: order.UpdatedAt.Format(time.RFC3339),
	}
}

func convertToModelCatalog(c *catalog.Catalog) *model.Catalog {
	var phoneModels []*model.Phone
	for _, phone := range c.Phones {
		phoneModels = append(phoneModels, convertToModelPhone(&phone))
	}

	return &model.Catalog{
		ID:        strconv.FormatUint(uint64(c.ID), 10),
		Name:      c.Name,
		Phones:    phoneModels,
		CreatedAt: c.CreatedAt.Format(time.RFC3339),
		UpdatedAt: c.UpdatedAt.Format(time.RFC3339),
	}
}

func convertToModelPhone(p *general.Phone) *model.Phone {
	return &model.Phone{
		ID:          strconv.FormatUint(uint64(p.ID), 10),
		Brand:       p.Brand,
		Model:       p.Name,
		Price:       p.Price,
		Description: p.Description,
	}
}
