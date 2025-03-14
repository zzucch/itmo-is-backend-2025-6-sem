package storage

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Storage struct {
	db *gorm.DB
}

func NewStorage(dsn string) (*Storage, error) {
	var err error

	db, err := gorm.Open(
		sqlite.Open(dsn),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		},
	)
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(
		&User{},
		&Phone{},
		&Image{},
		&Order{},
		&Catalog{},
		&Notification{},
	); err != nil {
		return nil, err
	}

	return &Storage{
		db: db,
	}, nil
}
