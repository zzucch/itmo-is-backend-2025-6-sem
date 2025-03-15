package server

import (
	"net/http"
)

func Setup(controllers controllers) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", controllers.indexController.Handle)
	mux.HandleFunc("/catalog", controllers.catalogController.Handle)
	mux.HandleFunc("/cart", controllers.cartController.Handle)
	mux.HandleFunc("/sell", controllers.sellController.Handle)
	mux.HandleFunc("/delete/{id}", controllers.sellController.HandleDelete)
	mux.HandleFunc("/create", controllers.sellController.HandleCreate)
	mux.HandleFunc("/phones", controllers.sellController.HandleGetAll)
	mux.HandleFunc("/notifications", controllers.notificationsHandler)
	mux.HandleFunc("/sse", controllers.notificationsSSEHandler)
	http.HandleFunc("/order", controllers.cartController.HandleOrder)

	mux.Handle(
		"/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))),
	)

	return mux
}
