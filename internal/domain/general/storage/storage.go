package storage

import (
	"errors"

	"github.com/is-web-y26/m3302-milovatskiy/internal/domain/cart"
	"github.com/is-web-y26/m3302-milovatskiy/internal/domain/catalog"
	"github.com/is-web-y26/m3302-milovatskiy/internal/domain/general"
	"github.com/is-web-y26/m3302-milovatskiy/internal/domain/user"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Storage struct {
	DB *gorm.DB
}

func (s *Storage) AddPhoneToCart(userID uint, phoneID uint) error {
	var user user.User
	if err := s.DB.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	var phone general.Phone
	if err := s.DB.First(&phone, phoneID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("phone not found")
		}
		return err
	}

	err := s.DB.Model(&user).Association("Cart").Append(&phone)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) RemovePhoneFromCart(userID uint, phoneID uint) error {
	var user user.User
	if err := s.DB.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	var phone general.Phone
	if err := s.DB.First(&phone, phoneID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("phone not found")
		}
		return err
	}

	err := s.DB.Model(&user).Association("Cart").Delete(&phone)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) FindCatalogsByUserID(userID uint) ([]catalog.Catalog, error) {
	if userID == 0 {
		return nil, errors.New("invalid user ID")
	}
	var catalogs []catalog.Catalog
	err := s.DB.Preload("Phones").Where("user_id = ?", userID).Find(&catalogs).Error
	if err != nil {
		return nil, err
	}
	return catalogs, nil
}

func (s *Storage) DeleteExpiredTokens() error {
	result := s.DB.Where("expires_at < ?", gorm.Expr("NOW()")).Delete(&user.Token{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *Storage) GetAllUsers() ([]*user.User, error) {
	var users []*user.User
	err := s.DB.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *Storage) UpdatePhone(phone *general.Phone) error {
	if phone == nil {
		return errors.New("phone cannot be nil")
	}
	if phone.ID == 0 {
		return errors.New("invalid phone ID")
	}

	var existingPhone general.Phone
	if err := s.DB.First(&existingPhone, phone.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("phone not found")
		}
		return err
	}

	return s.DB.Save(phone).Error
}

func (s *Storage) FindPhoneByID(id uint) (*general.Phone, error) {
	if id == 0 {
		return nil, errors.New("invalid ID")
	}
	var phone general.Phone
	err := s.DB.First(&phone, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &phone, nil
}

func (s *Storage) CreateUser(user *user.User) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}
	return s.DB.Create(user).Error
}

func (s *Storage) FindUserByID(id uint) (*user.User, error) {
	if id == 0 {
		return nil, errors.New("invalid ID")
	}
	var user user.User
	err := s.DB.Preload("Tokens").First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (s *Storage) FindUserByUsername(username string) (*user.User, error) {
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}
	var user user.User
	err := s.DB.Preload("Tokens").Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (s *Storage) UpdateUser(user *user.User) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}
	return s.DB.Save(user).Error
}

func (s *Storage) DeleteUser(id uint) error {
	if id == 0 {
		return errors.New("invalid ID")
	}
	result := s.DB.Delete(&user.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (s *Storage) CreateToken(token *user.Token) error {
	if token == nil {
		return errors.New("token cannot be nil")
	}

	if err := s.DB.Create(token).Error; err != nil {
		return err
	}

	return s.DB.Model(&user.User{Model: gorm.Model{ID: token.UserID}}).
		Association("Tokens").
		Append(token)
}

func (s *Storage) FindToken(tokenString string) (*user.Token, error) {
	if tokenString == "" {
		return nil, errors.New("token string cannot be empty")
	}
	var token user.Token
	err := s.DB.Where("token = ?", tokenString).First(&token).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &token, nil
}

func (s *Storage) DeleteToken(tokenString string) error {
	if tokenString == "" {
		return errors.New("token string cannot be empty")
	}
	result := s.DB.Where("token = ?", tokenString).Delete(&user.Token{})
	if result.Error != nil {
		return result.Error
	}
	return nil
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
		&cart.Order{},
		&catalog.Catalog{},
	); err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(
		&user.User{},
		&user.Token{},
		&cart.Order{},
		&general.Phone{},
		&general.Image{},
		&catalog.Catalog{},
	); err != nil {
		return nil, err
	}

	return &Storage{
		DB: db,
	}, nil
}
