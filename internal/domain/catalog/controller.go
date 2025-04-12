package catalog

import (
	"html/template"
	"net/http"

	pageData "github.com/is-web-y26/m3302-milovatskiy/internal/domain/general/page_data"
)

type Controller struct {
	service  *Service
	pageData *pageData.CatalogPageData
	template *template.Template
}

func NewController(
	service *Service,
	template *template.Template,
) *Controller {
	return &Controller{
		service:  service,
		pageData: pageData.NewCatalogPageData(),
		template: template,
	}
}

// @Summary Renders catalog page
// @Description Renders catalog page
// @Tags pages
// @Router /catalog [get]
// @Produce html
// @Success 200 "HTML catalog page"
// @Failure 500 "Internal server error"
func (c *Controller) Handle(
	responseWriter http.ResponseWriter,
	_ *http.Request,
) {
	var err error

	if c.pageData.SalePhone, err = c.service.GetSalePhone(); err != nil {
		http.Error(
			responseWriter,
			err.Error(),
			http.StatusInternalServerError,
		)

		return
	}

	if c.pageData.NewPhones, err = c.service.GetNewPhones(); err != nil {
		http.Error(
			responseWriter,
			err.Error(),
			http.StatusInternalServerError,
		)

		return
	}

	if c.pageData.FeaturedPhones, err = c.service.GetFeaturedPhones(); err != nil {
		http.Error(
			responseWriter,
			err.Error(),
			http.StatusInternalServerError,
		)

		return
	}

	if err := c.template.ExecuteTemplate(
		responseWriter,
		"layout",
		c.pageData,
	); err != nil {
		http.Error(
			responseWriter,
			err.Error(),
			http.StatusInternalServerError,
		)
	}
}
