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

// @Summary adds a new phone listing
// @Description adds a new phone listing
// @Tags phones
// @Router /api/phones [post]
// @Security BearerAuth
// @Param request body general.Phone true "Phone details"
// @Success 201
func (c *Controller) HandleCreatePhone(
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
			"invalid request body: "+err.Error(),
			http.StatusBadRequest,
		)

		return
	}

	userID, ok := request.Context().Value("user_id").(uint)
	if !ok {
		http.Error(responseWriter, "unauthorized", http.StatusUnauthorized)
		return
	}

	phone.SellerID = userID

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

// @Summary lists user's phones
// @Description lists user's phones
// @Tags phones
// @Router /api/phones [get]
// @Security BearerAuth
// @Success 200 {array} general.Phone
func (c *Controller) HandleGetAll(
	responseWriter http.ResponseWriter,
	r *http.Request,
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

	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		http.Error(responseWriter, "unauthorized", http.StatusUnauthorized)
		return
	}

	isAdmin, ok := r.Context().Value("is_admin").(bool)
	if !ok {
		isAdmin = false
	}

	userPhones := make([]general.Phone, 0, len(phones))
	for _, phone := range phones {
		if phone.SellerID == userID {
			userPhones = append(userPhones, phone)
		} else if isAdmin {
			userPhones = append(userPhones, phone)
		}
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(responseWriter).Encode(userPhones)
}

// @Summary removes a phone listing
// @Description removes a phone listing
// @Tags phones
// @Router /api/phones/{id} [delete]
// @Security BearerAuth
// @Param id path int true "Phone ID"
// @Success 204
func (c *Controller) HandleDeletePhone(
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

	phoneIDStr := request.URL.Path[len("/api/phones/"):]
	phoneID, err := strconv.Atoi(phoneIDStr)
	if err != nil {
		http.Error(
			responseWriter,
			"invalid phone ID: "+err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	userID, ok := request.Context().Value("user_id").(uint)
	if !ok {
		http.Error(responseWriter, "unauthorized", http.StatusUnauthorized)
		return
	}

	isAdmin, ok := request.Context().Value("is_admin").(bool)
	if !ok {
		isAdmin = false
	}

	phone, err := c.service.repository.FindPhoneByID(uint(phoneID))
	if err != nil {
		http.Error(responseWriter, "phone does not exist", http.StatusBadRequest)
		return
	}

	if phone.SellerID != userID && !isAdmin {
		http.Error(responseWriter, "unauthorized", http.StatusUnauthorized)
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

// @Summary gets phone listing by id
// @Description gets phone listing by id
// @Tags phones
// @Router /api/phones/{id} [get]
// @Security BearerAuth
// @Param id path int true "Phone ID"
// @Param request body general.Phone true "Phone details"
// @Success 200 {object} general.Phone
func (c *Controller) HandleGetPhoneByID(
	responseWriter http.ResponseWriter,
	request *http.Request,
) {
	if request.Method != http.MethodGet {
		http.Error(
			responseWriter,
			"invalid request method",
			http.StatusMethodNotAllowed,
		)
		return
	}

	phoneIDStr := request.URL.Path[len("/api/phones/"):]
	phoneID, err := strconv.Atoi(phoneIDStr)
	if err != nil {
		http.Error(
			responseWriter,
			"invalid phone ID",
			http.StatusBadRequest,
		)
		return
	}

	userID, ok := request.Context().Value("user_id").(uint)
	if !ok {
		http.Error(responseWriter, "unauthorized", http.StatusUnauthorized)
		return
	}

	isAdmin, ok := request.Context().Value("is_admin").(bool)
	if !ok {
		isAdmin = false
	}

	phone, err := c.service.GetPhoneByID(uint(phoneID))
	if err != nil {
		http.Error(responseWriter, "phone not found", http.StatusNotFound)
		return
	}

	if phone.SellerID != userID && !isAdmin {
		http.Error(responseWriter, "unauthorized", http.StatusUnauthorized)
		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(responseWriter).Encode(phone)
}

// @Summary modifies phone listing
// @Description modifies phone listing
// @Tags phones
// @Router /api/phones/{id} [put]
// @Security BearerAuth
// @Param id path int true "Phone ID"
// @Param request body general.Phone true "Updated phone details"
// @Success 200 {object} general.Phone
func (c *Controller) HandleUpdatePhone(
	responseWriter http.ResponseWriter,
	request *http.Request,
) {
	if request.Method != http.MethodPut {
		http.Error(
			responseWriter,
			"invalid request method",
			http.StatusMethodNotAllowed,
		)
		return
	}

	phoneIDStr := request.URL.Path[len("/api/phones/"):]
	phoneID, err := strconv.Atoi(phoneIDStr)
	if err != nil {
		http.Error(
			responseWriter,
			"invalid phone ID",
			http.StatusBadRequest,
		)
		return
	}

	userID, ok := request.Context().Value("user_id").(uint)
	if !ok {
		http.Error(responseWriter, "unauthorized", http.StatusUnauthorized)
		return
	}

	isAdmin, ok := request.Context().Value("is_admin").(bool)
	if !ok {
		isAdmin = false
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

	existingPhone, err := c.service.GetPhoneByID(uint(phoneID))
	if err != nil {
		http.Error(responseWriter, "phone not found", http.StatusNotFound)
		return
	}

	if existingPhone.SellerID != userID && !isAdmin {
		http.Error(responseWriter, "unauthorized", http.StatusUnauthorized)
		return
	}

	phone.ID = uint(phoneID)
	phone.SellerID = userID

	if err := c.service.UpdatePhone(&phone); err != nil {
		http.Error(
			responseWriter,
			"failed to update phone: "+err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(responseWriter).Encode(phone)
}
