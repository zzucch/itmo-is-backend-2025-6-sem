package resolvers

import (
	"strconv"
	"time"

	"github.com/is-web-y26/m3302-milovatskiy/internal/domain/cart"
	"github.com/is-web-y26/m3302-milovatskiy/internal/domain/graph/model"
	"github.com/is-web-y26/m3302-milovatskiy/internal/domain/sell"
	"github.com/is-web-y26/m3302-milovatskiy/internal/domain/user"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	PhoneService *sell.Service
	UserService  *user.Service
	CartService  *cart.Service
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
