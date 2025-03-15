package catalog

import (
	"github.com/is-web-y26/m3302-milovatskiy/internal/domain/general"
)

type CatalogService struct {
	repository CatalogRepository
}

func NewCatalogService(repository CatalogRepository) *CatalogService {
	return &CatalogService{repository: repository}
}

func (s *CatalogService) GetAllCatalogs() ([]Catalog, error) {
	return s.repository.FindAllCatalogs()
}

func (s *CatalogService) GetPhonesInCatalog(catalogID uint) (*Catalog, error) {
	return s.repository.FindCatalogByID(catalogID)
}

func (s *CatalogService) GetSalePhone() (general.Phone, error) {
	catalog, err := s.repository.FindCatalogByID(1)
	if err != nil {
		return general.Phone{}, err
	}

	return catalog.Phones[0], nil
}

func (s *CatalogService) GetFeaturedPhones() ([]general.Phone, error) {
	catalog, err := s.repository.FindCatalogByID(2)
	if err != nil {
		return nil, err
	}

	return catalog.Phones, nil
}

func (s *CatalogService) GetNewPhones() ([]general.Phone, error) {
	catalog, err := s.repository.FindCatalogByID(3)
	if err != nil {
		return nil, err
	}

	return catalog.Phones, nil
}
