package server

import (
	"net/http"

	"github.com/is-web-y26/m3302-milovatskiy/internal/handler"
)

func Setup() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handler.HandleIndex)
	mux.HandleFunc("/catalog", handler.HandleCatalog)
	mux.HandleFunc("/cart", handler.HandleCart)
	mux.HandleFunc("/sell", handler.HandleSell)
	mux.Handle(
		"/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))),
	)

	return mux
}
