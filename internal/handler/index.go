package handler

import (
	"net/http"
)

func HandleIndex(responseWriter http.ResponseWriter, request *http.Request) {
	HandleCommon(responseWriter, request, "templates/index.html")
}
