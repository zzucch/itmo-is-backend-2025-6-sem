package catalog

import (
	"errors"

	"github.com/is-web-y26/m3302-milovatskiy/internal/domain/general"
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) GetAllCatalogs() ([]Catalog, error) {
	return s.repository.FindAllCatalogs()
}

func (s *Service) GetPhonesInCatalog(catalogID uint) (*Catalog, error) {
	return s.repository.FindCatalogByID(catalogID)
}

func (s *Service) GetSalePhone() (general.Phone, error) {
	catalog, err := s.repository.FindCatalogByID(1)
	if err != nil {
		return general.Phone{}, err
	}

	if catalog == nil {
		return general.Phone{}, errors.New("does not exist")
	}

	if len(catalog.Phones) == 0 {
		return general.Phone{}, nil
	}

	return catalog.Phones[0], nil
}

func (s *Service) GetFeaturedPhones() ([]general.Phone, error) {
	catalog, err := s.repository.FindCatalogByID(2)
	if err != nil {
		return nil, err
	}

	return catalog.Phones, nil
}

func (s *Service) GetNewPhones() ([]general.Phone, error) {
	catalog, err := s.repository.FindCatalogByID(3)
	if err != nil {
		return nil, err
	}

	return catalog.Phones, nil
}

func (s *Service) FindAllCatalogs() ([]Catalog, error) {
	return s.repository.FindAllCatalogs()
}

func (s *Service) FindCatalogsByUserID(userID uint) ([]Catalog, error) {
	return s.repository.FindCatalogsByUserID(userID)
}

func (s *Service) GetCatalogByID(param any) (Catalog, error) {
	id, ok := param.(uint)
	if !ok {
		return Catalog{}, errors.New("invalid ID type")
	}

	catalog, err := s.repository.FindCatalogByID(id)
	if err != nil {
		return Catalog{}, err
	}

	if catalog == nil {
		return Catalog{}, errors.New("catalog does not exist")
	}

	return *catalog, nil
}

func (s *Service) CreateCatalog(catalog *Catalog) error {
	return s.repository.CreateCatalog(catalog)
}

func (s *Service) UpdateCatalog(catalog *Catalog) error {
	return s.repository.UpdateCatalog(catalog)
}

func (s *Service) DeleteCatalog(param any) error {
	id, ok := param.(uint)
	if !ok {
		return errors.New("invalid ID type")
	}
	return s.repository.DeleteCatalog(id)
}
