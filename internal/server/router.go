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
	// mux.HandleFunc("/login", handler.HandleLogin)
	// mux.HandleFunc("/signup", handler.HandleLogin)
	// mux.HandleFunc("/signout", handler.HandleLogout)

	mux.Handle(
		"/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))),
	)

	return mux
}
