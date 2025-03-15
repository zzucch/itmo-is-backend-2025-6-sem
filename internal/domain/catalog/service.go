package catalog

import (
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
