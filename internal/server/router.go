package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/is-web-y26/m3302-milovatskiy/docs"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

// @title Phone Marketplace API
// @version 1.0
// @BasePath /
func Setup(controllers controllers) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", controllers.indexController.Handle)
	mux.HandleFunc("/cart", controllers.cartController.Handle)
	mux.HandleFunc("/sell", controllers.sellController.Handle)
	mux.HandleFunc("/login", controllers.userController.HandleLoginPage)
	mux.HandleFunc("/signup", controllers.userController.HandleSignupPage)
	mux.HandleFunc("/order", controllers.cartController.HandleOrder)
	mux.HandleFunc("/catalog", controllers.catalogController.Handle)
	mux.HandleFunc("/notifications", controllers.notificationsHandler)

	go func() {
		gin.SetMode(gin.ReleaseMode)
		r := gin.New()

		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		// it conflicts with "/" on 3000
		r.Run(":3001")
	}()

	mux.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			controllers.userController.CreateUser(w, r)
			return
		}

		controllers.userController.AuthMiddleware(
			controllers.userController.AdminMiddleware(
				func(w http.ResponseWriter, r *http.Request) {
					switch r.Method {
					case http.MethodGet:
						controllers.userController.GetAllUsers(w, r)
					case http.MethodPut:
						controllers.userController.UpdateUser(w, r)
					case http.MethodDelete:
						controllers.userController.DeleteUser(w, r)
					default:
						http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
					}
				},
			),
		)
	})
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
			case http.MethodPost:
				controllers.sellController.HandleCreate(w, r)
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
	mux.HandleFunc("/api/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controllers.userController.GetUserByID(w, r)
		case http.MethodPut:
			controllers.userController.UpdateUser(w, r)
		case http.MethodDelete:
			controllers.userController.DeleteUser(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/sse", controllers.notificationsSSEHandler)

	mux.Handle(
		"/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))),
	)

	return mux
}
