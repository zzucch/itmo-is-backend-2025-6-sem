package handler

import (
	"net/http"
)

func HandleCatalog(responseWriter http.ResponseWriter, request *http.Request) {
	HandleCommon(responseWriter, request, "templates/catalog.html")
}
