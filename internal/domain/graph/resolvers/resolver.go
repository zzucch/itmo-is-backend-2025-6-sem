package resolvers

import (
	"github.com/is-web-y26/m3302-milovatskiy/internal/domain/sell"
	"github.com/is-web-y26/m3302-milovatskiy/internal/domain/user"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	PhoneService *sell.Service
	UserService  *user.Service
}
