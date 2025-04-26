package server

import (
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	_ "github.com/is-web-y26/m3302-milovatskiy/docs"
	"github.com/is-web-y26/m3302-milovatskiy/internal/domain/graph/generated"
	"github.com/is-web-y26/m3302-milovatskiy/internal/domain/graph/resolvers"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"github.com/vektah/gqlparser/v2/ast"
)

func Setup(controllers controllers, services services) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", controllers.indexController.Handle)
	mux.HandleFunc("/cart", controllers.cartController.Handle)
	mux.HandleFunc("/sell", controllers.sellController.Handle)
	mux.HandleFunc("/login", controllers.userController.HandleLoginPage)
	mux.HandleFunc("/signup", controllers.userController.HandleSignupPage)
	mux.HandleFunc("/catalog", controllers.catalogController.Handle)
	mux.HandleFunc("/notifications", controllers.notificationsHandler)

	c := generated.Config{
		Resolvers: &resolvers.Resolver{
			PhoneService: services.sellService,
			UserService:  services.userService,
		},
	}

	graphqlServer := handler.New(generated.NewExecutableSchema(c))

	graphqlServer.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
	})
	graphqlServer.AddTransport(transport.Options{})
	graphqlServer.AddTransport(transport.GET{})
	graphqlServer.AddTransport(transport.POST{})
	graphqlServer.AddTransport(transport.MultipartForm{})

	graphqlServer.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	graphqlServer.Use(extension.Introspection{})
	graphqlServer.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	const complexityLimit = 2000
	graphqlServer.Use(extension.FixedComplexityLimit(complexityLimit))

	mux.Handle("/graphql", playground.Handler(
		"The Phone Marketplace GraphQL playground",
		"/query",
	))
	mux.Handle("/query", graphqlServer)

	go func() {
		gin.SetMode(gin.ReleaseMode)
		r := gin.New()
		r.Use(func(c *gin.Context) {
			c.Writer.Header().Set(
				"Access-Control-Allow-Origin",
				"*",
			)
			c.Writer.Header().Set(
				"Access-Control-Allow-Methods",
				"GET, POST, PUT, DELETE, OPTIONS",
			)
			c.Writer.Header().Set(
				"Access-Control-Allow-Headers",
				"Content-Type, Authorization",
			)

			if c.Request.Method == http.MethodOptions {
				c.AbortWithStatus(http.StatusNoContent)
				return
			}

			c.Next()
		})

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
		)(w, r)
	})
	mux.HandleFunc("/api/users/login", controllers.userController.Login)
	mux.HandleFunc("/api/users/logout", controllers.userController.Logout)
	mux.HandleFunc("/api/users/me", controllers.userController.AuthMiddleware(
		controllers.userController.AdminMiddleware(
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
		),
	))
	mux.HandleFunc("/api/users/me/tokens", controllers.userController.AuthMiddleware(
		controllers.userController.AdminMiddleware(
			func(w http.ResponseWriter, r *http.Request) {
				switch r.Method {
				case http.MethodGet:
					controllers.userController.GetUserTokens(w, r)
				case http.MethodPost:
					controllers.sellController.HandleCreatePhone(w, r)
				default:
					http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				}
			},
		),
	))
	mux.HandleFunc("/api/users/{id}", controllers.userController.AuthMiddleware(
		controllers.userController.AdminMiddleware(
			func(w http.ResponseWriter, r *http.Request) {
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
			},
		),
	))
	mux.HandleFunc("/api/users/me/cart", controllers.userController.AuthMiddleware(
		func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				controllers.userController.GetCart(w, r)
			case http.MethodPost:
				controllers.userController.AddToCart(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		},
	))
	mux.HandleFunc("/api/users/me/cart/{phone_id}", controllers.userController.AuthMiddleware(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				controllers.userController.RemoveFromCart(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		},
	))
	mux.HandleFunc("/api/phones", controllers.userController.AuthMiddleware(
		controllers.userController.AdminMiddleware(
			func(w http.ResponseWriter, r *http.Request) {
				switch r.Method {
				case http.MethodGet:
					controllers.sellController.HandleGetAll(w, r)
				case http.MethodPost:
					controllers.sellController.HandleCreatePhone(w, r)
				default:
					http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				}
			},
		),
	))
	mux.HandleFunc("/api/phones/{id}", controllers.userController.AuthMiddleware(
		controllers.userController.AdminMiddleware(
			func(w http.ResponseWriter, r *http.Request) {
				switch r.Method {
				case http.MethodDelete:
					controllers.sellController.HandleDeletePhone(w, r)
				case http.MethodPost:
					controllers.sellController.HandleCreatePhone(w, r)
				case http.MethodPut:
					controllers.sellController.HandleUpdatePhone(w, r)
				case http.MethodGet:
					controllers.sellController.HandleGetPhoneByID(w, r)
				default:
					http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				}
			},
		),
	))
	mux.HandleFunc("/api/orders", controllers.userController.AuthMiddleware(
		controllers.userController.AdminMiddleware(
			func(w http.ResponseWriter, r *http.Request) {
				switch r.Method {
				case http.MethodGet:
					controllers.cartController.GetAllOrders(w, r)
				case http.MethodPost:
					controllers.cartController.CreateOrder(w, r)
				default:
					http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				}
			},
		),
	))
	mux.HandleFunc("/api/orders/{id}", controllers.userController.AuthMiddleware(
		controllers.userController.AdminMiddleware(
			func(w http.ResponseWriter, r *http.Request) {
				switch r.Method {
				case http.MethodDelete:
					controllers.cartController.DeleteOrder(w, r)
				case http.MethodPut:
					controllers.cartController.UpdateOrder(w, r)
				case http.MethodGet:
					controllers.cartController.GetOrderByID(w, r)
				default:
					http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				}
			},
		),
	))
	mux.HandleFunc("/api/orders/me", controllers.userController.AuthMiddleware(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet {
				controllers.cartController.GetMyOrders(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		},
	))
	mux.HandleFunc("/api/catalogs", controllers.userController.AuthMiddleware(
		controllers.userController.AdminMiddleware(
			func(w http.ResponseWriter, r *http.Request) {
				switch r.Method {
				case http.MethodGet:
					controllers.catalogController.HandleGetCatalogs(w, r)
				case http.MethodPost:
					controllers.catalogController.HandleCreateCatalog(w, r)
				default:
					http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				}
			},
		),
	))
	mux.HandleFunc("/api/catalogs/{id}", controllers.userController.AuthMiddleware(
		controllers.userController.AdminMiddleware(
			func(w http.ResponseWriter, r *http.Request) {
				switch r.Method {
				case http.MethodDelete:
					controllers.catalogController.HandleDeleteCatalog(w, r)
				case http.MethodPut:
					controllers.catalogController.HandleUpdateCatalog(w, r)
				case http.MethodGet:
					controllers.catalogController.HandleGetCatalogByID(w, r)
				default:
					http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				}
			},
		),
	))

	mux.HandleFunc("/api/sse", controllers.notificationsSSEHandler)

	mux.Handle(
		"/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))),
	)

	return timeMiddleware(
		cacheControlMiddleware(
			corsMiddleware(
				mux,
			),
		),
	)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
