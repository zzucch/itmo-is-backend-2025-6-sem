package catalog

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"

	pageData "github.com/is-web-y26/m3302-milovatskiy/internal/domain/general/pagedata"
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

// @Summary create catalog
// @Description Creates a new catalog entry for the logged-in user.
// @Tags catalogs
// @Accept json
// @Produce json
// @Param request body catalog.Catalog true "Catalog object to be created"
// @Success 201 {object} catalog.Catalog "Created catalog object"
// @Failure 400 {string} string "Invalid request body"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Could not create catalog"
// @Router /api/catalogs [post]
// @Security BearerAuth
func (c *Controller) HandleCreateCatalog(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var catalog Catalog
	if err := json.NewDecoder(r.Body).Decode(&catalog); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	catalog.CreatorID = userID

	if err := c.service.CreateCatalog(&catalog); err != nil {
		http.Error(w, "could not create catalog", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(catalog); err != nil {
		log.Fatal(err)
	}
}

// @Summary get all catalogs
// @Description Retrieves all catalogs for admin or user's own catalogs.
// @Tags catalogs
// @Produce json
// @Success 200 {array} catalog.Catalog "List of catalogs"
// @Failure 500 {string} string "Could not fetch catalogs"
// @Router /api/catalogs [get]
// @Security BearerAuth
func (c *Controller) HandleGetCatalogs(w http.ResponseWriter, _ *http.Request) {
	var catalogs []Catalog
	var err error

	catalogs, err = c.service.FindAllCatalogs()
	if err != nil {
		http.Error(w, "could not fetch catalogs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(catalogs); err != nil {
		log.Fatal(err)
	}
}

// @Summary get catalog by ID
// @Description Retrieves a catalog by its ID. Only accessible by its creator or an admin.
// @Tags catalogs
// @Produce json
// @Param id path int true "Catalog ID"
// @Success 200 {object} catalog.Catalog "Catalog object"
// @Failure 400 {string} string "Invalid catalog ID"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Catalog not found"
// @Router /api/catalogs/{id} [get]
// @Security BearerAuth
func (c *Controller) HandleGetCatalogByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/api/catalogs/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid catalog ID", http.StatusBadRequest)
		return
	}

	catalog, err := c.service.GetCatalogByID(uint(id))
	if err != nil {
		http.Error(w, "catalog not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(catalog); err != nil {
		log.Fatal(err)
	}
}

// @Summary update catalog
// @Description Updates an existing catalog entry. Only accessible by its creator or an admin.
// @Tags catalogs
// @Accept json
// @Produce json
// @Param id path int true "Catalog ID"
// @Param request body catalog.Catalog true "Updated catalog data"
// @Success 200 {object} catalog.Catalog "Updated catalog object"
// @Failure 400 {string} string "Bad request or invalid catalog ID"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Catalog not found"
// @Failure 500 {string} string "Update failed"
// @Router /api/catalogs/{id} [put]
// @Security BearerAuth
func (c *Controller) HandleUpdateCatalog(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/api/catalogs/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid catalog ID", http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		log.Fatal("not ok")
	}
	isAdmin, ok := r.Context().Value("is_admin").(bool)
	if !ok {
		log.Fatal("not ok")
	}

	var updated Catalog
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	existing, err := c.service.GetCatalogByID(uint(id))
	if err != nil {
		http.Error(w, "catalog not found", http.StatusNotFound)
		return
	}

	if existing.CreatorID != userID && !isAdmin {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	updated.ID = uint(id)
	updated.CreatorID = existing.CreatorID

	if err := c.service.UpdateCatalog(&updated); err != nil {
		http.Error(w, "update failed", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(updated); err != nil {
		log.Fatal(err)
	}
}

// @Summary delete catalog
// @Description Deletes a catalog entry by its ID. Only accessible by its creator or an admin.
// @Tags catalogs
// @Param id path int true "Catalog ID"
// @Success 204 "No Content"
// @Failure 400 {string} string "Invalid catalog ID"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Catalog not found"
// @Failure 500 {string} string "Delete failed"
// @Router /api/catalogs/{id} [delete]
// @Security BearerAuth
func (c *Controller) HandleDeleteCatalog(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/api/catalogs/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid catalog ID", http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		log.Fatal("not ok")
	}
	isAdmin, ok := r.Context().Value("is_admin").(bool)
	if !ok {
		log.Fatal("not ok")
	}

	catalog, err := c.service.GetCatalogByID(uint(id))
	if err != nil {
		http.Error(w, "catalog not found", http.StatusNotFound)
		return
	}

	if catalog.CreatorID != userID && !isAdmin {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if err := c.service.DeleteCatalog(uint(id)); err != nil {
		http.Error(w, "delete failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
