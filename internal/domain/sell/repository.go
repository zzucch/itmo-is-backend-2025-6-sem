package sell

import "github.com/is-web-y26/m3302-milovatskiy/internal/domain/general"

type Repository interface {
	CreatePhone(phone *general.Phone) error
	FindAllPhones() ([]general.Phone, error)
	FindPhoneByID(id uint) (*general.Phone, error)
	DeletePhone(id uint) error
	UpdatePhone(phone *general.Phone) error
}
