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

// HandleCreate adds a new phone listing
// @Router /api/phones [post]
// @Security BearerAuth
// @Param request body general.Phone true "Phone details"
// @Success 201
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

// HandleGetAll lists user's phones
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

	userPhones := make([]general.Phone, 0, len(phones))
	for _, phone := range phones {
		if phone.SellerID == userID {
			userPhones = append(userPhones, phone)
		}
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(responseWriter).Encode(userPhones)
}

// HandleDelete removes a phone listing
// @Router /api/phones/{id} [delete]
// @Security BearerAuth
// @Param id path int true "Phone ID"
// @Success 204
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

	userID, ok := request.Context().Value("user_id").(uint)
	if !ok {
		http.Error(responseWriter, "unauthorized", http.StatusUnauthorized)
		return
	}

	phone, err := c.service.repository.FindPhoneByID(uint(phoneID))
	if err != nil {
		http.Error(responseWriter, "phone does not exist", http.StatusBadRequest)
		return
	}

	if phone.SellerID != userID {
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

// HandleUpdate modifies phone listing
// @Router /api/phones/{id} [put]
// @Security BearerAuth
// @Param id path int true "Phone ID"
// @Param request body general.Phone true "Updated phone details"
// @Success 200 {object} general.Phone
func (c *Controller) HandleUpdate(
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

	if existingPhone.SellerID != userID {
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
