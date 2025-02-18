package handler

import (
	"net/http"
)

func HandleSell(responseWriter http.ResponseWriter, request *http.Request) {
	HandleCommon(responseWriter, request, "templates/sell.html")
}
