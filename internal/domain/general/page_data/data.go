package page_data

import "github.com/is-web-y26/m3302-milovatskiy/internal/domain/general"

type IndexPageData struct {
	*CommonPageData
	SalePhone      general.Phone
	NewPhones      []general.Phone
	FeaturedPhones []general.Phone
}

func NewIndexPageData() *IndexPageData {
	return &IndexPageData{
		CommonPageData: &DefaultData,
	}
}

type CatalogPageData struct {
	*CommonPageData
	SalePhone      general.Phone
	NewPhones      []general.Phone
	FeaturedPhones []general.Phone
}

func NewCatalogPageData() *CatalogPageData {
	return &CatalogPageData{
		CommonPageData: &DefaultData,
	}
}

type CartPageData struct {
	*CommonPageData
	CartItems []general.Phone
}

func NewCartPageData() *CartPageData {
	return &CartPageData{
		CommonPageData: &DefaultData,
	}
}

type SellPageData struct {
	*CommonPageData
	SellPhones []general.Phone
}

func NewSellPageData() *SellPageData {
	return &SellPageData{
		CommonPageData: &DefaultData,
	}
}

type CommonPageData struct {
	Title        string
	MenuLinks    []MenuLink
	IsAuthorized bool
}

type MenuLink struct {
	Name string
	URL  string
}

var DefaultData = CommonPageData{
	Title: "The Phone Marketplace",
	MenuLinks: []MenuLink{
		{Name: "Home", URL: "/"},
		{Name: "Catalog", URL: "/catalog"},
		{Name: "Sell", URL: "/sell"},
		{Name: "Contacts", URL: "#contacts"},
		{Name: "Cart", URL: "/cart"},
		{Name: "notifications", URL: "/notifications"},
	},
	IsAuthorized: false,
}
