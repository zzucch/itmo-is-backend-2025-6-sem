package catalog

type Repository interface {
	FindAllCatalogs() ([]Catalog, error)
	FindCatalogByID(id uint) (*Catalog, error)
	CreateCatalog(catalog *Catalog) error
	UpdateCatalog(catalog *Catalog) error
	DeleteCatalog(id uint) error
}
