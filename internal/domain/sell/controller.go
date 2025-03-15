package sell

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"

	"github.com/is-web-y26/m3302-milovatskiy/internal/domain/general"
	pageData "github.com/is-web-y26/m3302-milovatskiy/internal/domain/general/page_data"
)

type Controller struct {
	service  *Service
	pageData *pageData.SellPageData
	template *template.Template
}

func NewController(
	service *Service,
	template *template.Template,
) *Controller {
	return &Controller{
		service:  service,
		pageData: pageData.NewSellPageData(),
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

func (c *Controller) HandleCreate(
	responseWriter http.ResponseWriter,
	request *http.Request,
) {
	if request.Method != http.MethodPost {
		http.Error(
			responseWriter,
			"invalid request method",
			http.StatusMethodNotAllowed,
		)
		return
	}

	var phone general.Phone
	if err := json.NewDecoder(request.Body).Decode(&phone); err != nil {
		http.Error(
			responseWriter,
			"invalid request body",
			http.StatusBadRequest,
		)

		return
	}

	if err := c.service.CreatePhone(&phone); err != nil {
		http.Error(
			responseWriter,
			"failed to create phone"+err.Error(),
			http.StatusInternalServerError,
		)

		return
	}

	responseWriter.WriteHeader(http.StatusCreated)
}

func (c *Controller) HandleGetAll(
	responseWriter http.ResponseWriter,
	_ *http.Request,
) {
	phones, err := c.service.FindAllPhones()
	if err != nil {
		http.Error(
			responseWriter,
			"failed to fetch phones",
			http.StatusInternalServerError,
		)

		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(responseWriter).Encode(phones)
}

func (c *Controller) HandleDelete(
	responseWriter http.ResponseWriter,
	request *http.Request,
) {
	if request.Method != http.MethodDelete {
		http.Error(
			responseWriter,
			"invalid request method",
			http.StatusMethodNotAllowed,
		)

		return
	}

	phoneIDStr := request.URL.Path[len("/delete/"):]
	phoneID, err := strconv.Atoi(phoneIDStr)
	if err != nil {
		http.Error(
			responseWriter,
			"invalid phone ID",
			http.StatusBadRequest,
		)
		return
	}

	if err := c.service.DeletePhone(uint(phoneID)); err != nil {
		http.Error(
			responseWriter,
			"failed to delete phone: "+err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	responseWriter.WriteHeader(http.StatusNoContent)
}
