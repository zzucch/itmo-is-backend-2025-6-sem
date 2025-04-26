package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/is-web-y26/m3302-milovatskiy/internal/config"
	"github.com/is-web-y26/m3302-milovatskiy/internal/domain/cart"
	"github.com/is-web-y26/m3302-milovatskiy/internal/domain/catalog"
	"github.com/is-web-y26/m3302-milovatskiy/internal/domain/general/index"
	"github.com/is-web-y26/m3302-milovatskiy/internal/domain/general/storage"
	"github.com/is-web-y26/m3302-milovatskiy/internal/domain/general/storage/s3"
	"github.com/is-web-y26/m3302-milovatskiy/internal/domain/sell"
	"github.com/is-web-y26/m3302-milovatskiy/internal/domain/user"
)

func Start(logger *log.Logger, config *config.Config) error {
	address := fmt.Sprintf(":%s", config.Port)

	s3Client, err := s3.NewS3Client()
	if err != nil {
		log.Fatalf("Failed to create S3 client: %v", err)
	}

	storage, err := storage.New(config.DSN)
	if err != nil {
		logger.Fatal(err)
	}

	storage.AddData()

	controllers, services := initControllers(logger, storage, s3Client)

	logger.Printf("server is running on %s", address)
	http.ListenAndServe(address, Setup(controllers, services))

	return nil
}

type controllers struct {
	indexController   *index.Controller
	catalogController *catalog.Controller
	sellController    *sell.Controller
	cartController    *cart.Controller
	userController    *user.Controller
}

type services struct {
	catalogService *catalog.Service
	sellService    *sell.Service
	cartService    *cart.Service
	userService    *user.Service
}

func initControllers(
	logger *log.Logger,
	storage *storage.Storage,
	s3Client *s3.S3Client,
) (controllers, services) {
	commonTemplates := []string{
		"templates/layout.html",
		"templates/header.html",
		"templates/footer.html",
	}

	indexTemplateFiles := append([]string{"templates/index.html"}, commonTemplates...)
	indexTemplate, err := template.ParseFiles(indexTemplateFiles...)
	if err != nil {
		logger.Fatal(err)
	}

	catalogTemplateFiles := append([]string{"templates/catalog.html"}, commonTemplates...)
	catalogTemplate, err := template.ParseFiles(catalogTemplateFiles...)
	if err != nil {
		logger.Fatal(err)
	}

	sellTemplateFiles := append([]string{"templates/sell.html"}, commonTemplates...)
	sellTemplate, err := template.ParseFiles(sellTemplateFiles...)
	if err != nil {
		logger.Fatal(err)
	}

	cartTemplateFiles := append([]string{"templates/cart.html"}, commonTemplates...)
	cartTemplate, err := template.ParseFiles(cartTemplateFiles...)
	if err != nil {
		logger.Fatal(err)
	}

	loginTemplateFiles := append([]string{"templates/login.html"}, commonTemplates...)
	loginTemplate, err := template.ParseFiles(loginTemplateFiles...)
	if err != nil {
		logger.Fatal(err)
	}

	signupTemplateFiles := append([]string{"templates/signup.html"}, commonTemplates...)
	signupTemplate, err := template.ParseFiles(signupTemplateFiles...)
	if err != nil {
		logger.Fatal(err)
	}

	sellService := sell.NewService(storage, s3Client)
	catalogService := catalog.NewService(storage)
	cartService := cart.NewService(storage)
	userService := user.NewService(storage)

	return controllers{
			index.NewController(catalogService, indexTemplate),
			catalog.NewController(catalogService, catalogTemplate),
			sell.NewController(sellService, sellTemplate),
			cart.NewController(cartService, cartTemplate),
			user.NewController(*userService, loginTemplate, signupTemplate),
		}, services{
			catalogService: catalogService,
			sellService:    sellService,
			cartService:    cartService,
			userService:    userService,
		}
}
