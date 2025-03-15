package cart

import (
	"html/template"
	"net/http"
	"strconv"

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

func (c *Controller) HandleOrder(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(responseWriter, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	userIDStr := request.FormValue("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(responseWriter, "Invalid user ID", http.StatusBadRequest)
		return
	}

	request.ParseForm()
	phoneIDStrings := request.Form["phone_ids"]

	var phoneIDs []uint
	for _, idStr := range phoneIDStrings {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(responseWriter, "Invalid phone ID", http.StatusBadRequest)
			return
		}
		phoneIDs = append(phoneIDs, uint(id))
	}

	order, err := c.service.PlaceOrder(uint(userID), phoneIDs)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	responseWriter.Write(
		[]byte("order ID: " + strconv.Itoa(int(order.ID))),
	)
}
