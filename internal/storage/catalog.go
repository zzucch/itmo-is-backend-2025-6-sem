package storage

import "gorm.io/gorm"

type Catalog struct {
	gorm.Model
	Name     string
	PhoneIDs []uint64
}
