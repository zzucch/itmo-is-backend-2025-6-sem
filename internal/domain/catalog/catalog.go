package catalog

import (
	"github.com/is-web-y26/m3302-milovatskiy/internal/domain/general"
	"gorm.io/gorm"
)

type Catalog struct {
	gorm.Model
	Name   string
	Phones []general.Phone `gorm:"many2many:catalog_phones;"`
}
