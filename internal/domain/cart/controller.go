package cart

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"
	"strings"

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
// @Router /api/orders [post]
// @Accept multipart/form-data
// @Param user_id formData int true "User ID"
// @Param phone_ids formData []int true "Array of phone IDs" collectionFormat: multi
// @Success 200 {string} string "order ID: {id}"
// @Failure 400 {string} string "Invalid user ID or phone ID"
// @Failure 405 {string} string "Method not allowed"
// @Failure 500 {string} string "Internal server error"
func (c *Controller) CreateOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	shippingAddress := r.PostFormValue("shipping_address")
	if shippingAddress == "" {
		http.Error(w, "Shipping address is required", http.StatusBadRequest)
		return
	}

	phoneIDs := make([]uint, 0)
	for _, idStr := range r.PostForm["phone_ids"] {
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			http.Error(w, "Invalid phone ID: "+idStr, http.StatusBadRequest)
			return
		}
		phoneIDs = append(phoneIDs, uint(id))
	}

	if len(phoneIDs) == 0 {
		http.Error(w, "No phones selected", http.StatusBadRequest)
		return
	}

	order, err := c.service.PlaceOrder(uint(userID), phoneIDs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"order_id": order.ID,
		"message":  "Order placed successfully",
	})
}

// @Summary Retrieves all orders
// @Description Retrieves all orders (admin only)
// @Tags orders
// @Router /api/orders [get]
// @Security BearerAuth
// @Success 200 {array} Order
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 500 {string} string "Internal server error"
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
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

// @Summary Retrieves an order by ID
// @Description Retrieves an order by ID (admin or order owner only)
// @Tags orders
// @Router /api/orders/{id} [get]
// @Param id path int true "Order ID"
// @Security BearerAuth
// @Success 200 {object} Order
// @Failure 400 {string} string "Invalid order ID"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {string} string "Order not found"
// @Failure 500 {string} string "Internal server error"
func (c *Controller) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/orders/")
	idStr := strings.Split(path, "/")[0]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	isAdmin, _ := r.Context().Value("is_admin").(bool)

	order, err := c.service.GetOrderByID(uint(id))
	if err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	if !isAdmin && order.UserID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(order); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

// @Summary Updates an order
// @Description Updates an order (admin only)
// @Tags orders
// @Router /api/orders/{id} [put]
// @Param id path int true "Order ID"
// @Security BearerAuth
// @Param request body Order true "Order data"
// @Success 200 {object} Order
// @Failure 400 {string} string "Invalid request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 500 {string} string "Internal server error"
func (c *Controller) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/orders/")
	idStr := strings.Split(path, "/")[0]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	if isAdmin, ok := r.Context().Value("is_admin").(bool); !isAdmin || !ok {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var order Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if uint(id) != order.ID {
		http.Error(w, "ID in path and body don't match", http.StatusBadRequest)
		return
	}

	updatedOrder, err := c.service.UpdateOrder(&order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updatedOrder); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

// @Summary Deletes an order
// @Description Deletes an order (admin only)
// @Tags orders
// @Router /api/orders/{id} [delete]
// @Param id path int true "Order ID"
// @Security BearerAuth
// @Success 204 "No Content"
// @Failure 400 {string} string "Invalid order ID"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 500 {string} string "Internal server error"
func (c *Controller) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/orders/")
	idStr := strings.Split(path, "/")[0]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	if isAdmin, ok := r.Context().Value("is_admin").(bool); !isAdmin || !ok {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if err := c.service.DeleteOrder(uint(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Summary Retrieves current user's orders
// @Description Retrieves orders for the authenticated user
// @Tags orders/me
// @Router /api/orders/me [get]
// @Security BearerAuth
// @Success 200 {array} Order
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
func (c *Controller) GetMyOrders(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	order, err := c.service.GetOrdersByUserID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(order); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
