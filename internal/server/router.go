package server

import (
	"net/http"
)

func Setup(controllers controllers) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", controllers.indexController.Handle)
	mux.HandleFunc("/cart", controllers.cartController.Handle)
	mux.HandleFunc("/sell", controllers.sellController.Handle)
	mux.HandleFunc("/login", controllers.userController.HandleLoginPage)
	mux.HandleFunc("/signup", controllers.userController.HandleSignupPage)
	http.HandleFunc("/order", controllers.cartController.HandleOrder)
	mux.HandleFunc("/catalog", controllers.catalogController.Handle)
	mux.HandleFunc("/notifications", controllers.notificationsHandler)

	mux.HandleFunc("/api/users", controllers.userController.CreateUser)
	mux.HandleFunc("/api/users/login", controllers.userController.Login)
	mux.HandleFunc("/api/users/logout", controllers.userController.Logout)

	mux.HandleFunc("/api/phones", controllers.userController.AuthMiddleware(
		func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				controllers.sellController.HandleGetAll(w, r)
			case http.MethodPost:
				controllers.sellController.HandleCreate(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		},
	))
	mux.HandleFunc("/api/phones/{id}", controllers.userController.AuthMiddleware(
		func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodDelete:
				controllers.sellController.HandleDelete(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		},
	))

	mux.HandleFunc("/api/users/me", controllers.userController.AuthMiddleware(
		func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				controllers.userController.GetCurrentUser(w, r)
			case http.MethodPut:
				controllers.userController.UpdateCurrentUser(w, r)
			case http.MethodDelete:
				controllers.userController.DeleteCurrentUser(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		},
	))

	mux.HandleFunc("/api/sse", controllers.notificationsSSEHandler)

	mux.Handle(
		"/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))),
	)

	return mux
}
