package handler

import "net/http"

func HandleLogin(responseWriter http.ResponseWriter, request *http.Request) {
	http.Redirect(responseWriter, request, "/?auth=true", http.StatusSeeOther)
}

func HandleLogout(responseWriter http.ResponseWriter, request *http.Request) {
	http.Redirect(responseWriter, request, "/?auth=false", http.StatusSeeOther)
}
