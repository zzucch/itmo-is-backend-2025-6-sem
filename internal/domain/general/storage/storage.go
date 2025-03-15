package storage

import (
	"errors"
	"os/user"

	"github.com/is-web-y26/m3302-milovatskiy/internal/domain/catalog"
	"github.com/is-web-y26/m3302-milovatskiy/internal/domain/general"
	"github.com/is-web-y26/m3302-milovatskiy/internal/domain/notification"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Storage struct {
	DB *gorm.DB
}

func (s *Storage) DeletePhone(id uint) error {
	if id == 0 {
		return errors.New("invalid ID")
	}

	result := s.DB.Delete(&general.Phone{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("phone not found")
	}

	return nil
}

func (s *Storage) FindAllPhones() ([]general.Phone, error) {
	var phones []general.Phone
	err := s.DB.Find(&phones).Error
	if err != nil {
		return nil, err
	}
	return phones, nil
}

func (s *Storage) CreatePhone(phone *general.Phone) error {
	if phone == nil {
		return errors.New("phone cannot be nil")
	}
	return s.DB.Create(phone).Error
}

func (s Storage) CreateCatalog(catalog *catalog.Catalog) error {
	if catalog == nil {
		return errors.New("catalog cannot be nil")
	}
	return s.DB.Create(catalog).Error
}

func (s Storage) DeleteCatalog(id uint) error {
	if id == 0 {
		return errors.New("invalid ID")
	}
	return s.DB.Delete(&catalog.Catalog{}, id).Error
}

func (s Storage) FindAllCatalogs() ([]catalog.Catalog, error) {
	var catalogs []catalog.Catalog
	err := s.DB.Preload("Phones").Find(&catalogs).Error
	if err != nil {
		return nil, err
	}
	return catalogs, nil
}

func (s Storage) FindCatalogByID(id uint) (*catalog.Catalog, error) {
	if id == 0 {
		return nil, errors.New("invalid ID")
	}
	var catalog catalog.Catalog
	err := s.DB.Preload("Phones").First(&catalog, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &catalog, nil
}

func (s Storage) UpdateCatalog(catalog *catalog.Catalog) error {
	if catalog == nil {
		return errors.New("catalog cannot be nil")
	}

	return s.DB.Save(catalog).Error
}

func New(dsn string) (*Storage, error) {
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

	if err := db.Migrator().DropTable(
		&user.User{},
		&general.Phone{},
		&general.Image{},
		&Order{},
		&catalog.Catalog{},
		&notification.Notification{},
	); err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(
		&user.User{},
		&general.Phone{},
		&general.Image{},
		&Order{},
		&catalog.Catalog{},
		&notification.Notification{},
	); err != nil {
		return nil, err
	}

	return &Storage{
		DB: db,
	}, nil
}
