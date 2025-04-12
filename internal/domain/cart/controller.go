package cart

import (
	"encoding/json"
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

// @Summary Renders cart page
// @Description Renders cart page
// @Tags pages
// @Router /cart [get]
// @Produce html
// @Success 200 "HTML cart page"
// @Failure 500 "Internal server error"
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

// @Summary Places a new order
// @Description Places a new order
// @Tags orders
// @Router /orders [post]
// @Accept multipart/form-data
// @Param user_id formData int true "User ID"
// @Param phone_ids formData []int true "Array of phone IDs" collectionFormat: multi
// @Success 200 {string} string "order ID: {id}"
// @Failure 400 {string} string "Invalid user ID or phone ID"
// @Failure 405 {string} string "Method not allowed"
// @Failure 500 {string} string "Internal server error"
func (c *Controller) CreateOrder(responseWriter http.ResponseWriter, request *http.Request) {
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

func (c *Controller) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	if _, ok := r.Context().Value("user_id").(uint); !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if isAdmin, ok := r.Context().Value("is_admin").(bool); !isAdmin && !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	users, err := c.service.GetAllOrders()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
