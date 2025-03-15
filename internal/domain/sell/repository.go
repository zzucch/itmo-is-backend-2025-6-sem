package sell

import "github.com/is-web-y26/m3302-milovatskiy/internal/domain/general"

type SellRepository interface {
	CreatePhone(phone *general.Phone) error
}
