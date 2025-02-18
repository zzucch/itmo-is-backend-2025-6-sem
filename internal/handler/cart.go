package handler

import (
	"net/http"
)

func HandleCart(responseWriter http.ResponseWriter, request *http.Request) {
	HandleCommon(responseWriter, request, "templates/cart.html")
}
