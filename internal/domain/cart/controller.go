package cart

import (
	"html/template"
	"net/http"

	pageData "github.com/is-web-y26/m3302-milovatskiy/internal/domain/general/page_data"
)

type Controller struct {
	service  *Service
	pageData *pageData.CartPageData
	template *template.Template
}

func NewController(
	service *Service,
	template *template.Template,
) *Controller {
	return &Controller{
		service:  service,
		pageData: pageData.NewCartPageData(),
		template: template,
	}
}

func (c *Controller) Handle(
	responseWriter http.ResponseWriter,
	_ *http.Request,
) {
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
