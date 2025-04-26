package catalog

import (
	"errors"
	"strconv"
	"time"

	"github.com/is-web-y26/m3302-milovatskiy/internal/domain/general"
)

type Service struct {
	repository Repository
	cache      *cache
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
		cache:      newCache(5 * time.Second),
	}
}

func (s *Service) GetAllCatalogs() ([]Catalog, error) {
	cacheKey := "all_catalogs"
	if cached, found := s.cache.get(cacheKey); found {
		return cached.([]Catalog), nil
	}

	catalogs, err := s.repository.FindAllCatalogs()
	if err != nil {
		return nil, err
	}

	s.cache.set(cacheKey, catalogs)
	return catalogs, nil
}

func (s *Service) GetPhonesInCatalog(catalogID uint) (*Catalog, error) {
	cacheKey := "catalog_" + strconv.FormatUint(uint64(catalogID), 10)
	if cached, found := s.cache.get(cacheKey); found {
		return cached.(*Catalog), nil
	}

	catalog, err := s.repository.FindCatalogByID(catalogID)
	if err != nil {
		return nil, err
	}

	s.cache.set(cacheKey, catalog)
	return catalog, nil
}

func (s *Service) GetSalePhone() (general.Phone, error) {
	cacheKey := "sale_phone"
	if cached, found := s.cache.get(cacheKey); found {
		return cached.(general.Phone), nil
	}

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

	phone := catalog.Phones[0]
	s.cache.set(cacheKey, phone)
	return phone, nil
}

func (s *Service) GetFeaturedPhones() ([]general.Phone, error) {
	cacheKey := "featured_phones"
	if cached, found := s.cache.get(cacheKey); found {
		return cached.([]general.Phone), nil
	}

	catalog, err := s.repository.FindCatalogByID(2)
	if err != nil {
		return nil, err
	}

	s.cache.set(cacheKey, catalog.Phones)
	return catalog.Phones, nil
}

func (s *Service) GetNewPhones() ([]general.Phone, error) {
	cacheKey := "new_phones"
	if cached, found := s.cache.get(cacheKey); found {
		return cached.([]general.Phone), nil
	}

	catalog, err := s.repository.FindCatalogByID(3)
	if err != nil {
		return nil, err
	}

	s.cache.set(cacheKey, catalog.Phones)
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
	err := s.repository.CreateCatalog(catalog)
	if err == nil {
		s.cache.invalidate("all_catalogs")
	}
	return err
}

func (s *Service) UpdateCatalog(catalog *Catalog) error {
	err := s.repository.UpdateCatalog(catalog)
	if err == nil {
		s.cache.invalidate("all_catalogs")
		s.cache.invalidate("catalog_" + strconv.FormatUint(uint64(catalog.ID), 10))
	}
	return err
}

func (s *Service) DeleteCatalog(param any) error {
	id, ok := param.(uint)
	if !ok {
		return errors.New("invalid ID type")
	}

	err := s.repository.DeleteCatalog(id)
	if err == nil {
		s.cache.invalidate("all_catalogs")
		s.cache.invalidate("catalog_" + strconv.FormatUint(uint64(id), 10))
	}
	return err
}
