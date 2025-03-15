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
	"github.com/is-web-y26/m3302-milovatskiy/internal/domain/sell"
)

func Start(logger *log.Logger, config *config.Config) error {
	address := fmt.Sprintf(":%s", config.Port)

	storage, err := storage.New(config.DSN)
	if err != nil {
		logger.Fatal(err)
	}

	storage.AddData()

	controllers := initControllers(logger, storage)

	logger.Printf("server is running on %s", address)
	http.ListenAndServe(address, Setup(controllers))

	return nil
}

type controllers struct {
	indexController   *index.Controller
	catalogController *catalog.Controller
	sellController    *sell.Controller
	cartController    *cart.Controller
}

func initControllers(
	logger *log.Logger,
	storage *storage.Storage,
) controllers {
	indexTemplate, err := template.ParseFiles(
		"templates/layout.html",
		"templates/header.html",
		"templates/footer.html",
		"templates/index.html",
	)
	if err != nil {
		logger.Fatal(err)
	}

	catalogTemplate, err := template.ParseFiles(
		"templates/layout.html",
		"templates/header.html",
		"templates/footer.html",
		"templates/catalog.html",
	)
	if err != nil {
		logger.Fatal(err)
	}

	sellTemplate, err := template.ParseFiles(
		"templates/layout.html",
		"templates/header.html",
		"templates/footer.html",
		"templates/sell.html",
	)
	if err != nil {
		logger.Fatal(err)
	}

	cartTemplate, err := template.ParseFiles(
		"templates/layout.html",
		"templates/header.html",
		"templates/footer.html",
		"templates/cart.html",
	)
	if err != nil {
		logger.Fatal(err)
	}

	sellService := sell.NewService(storage)
	catalogService := catalog.NewCatalogService(storage)
	cartService := &cart.CartService{}

	return controllers{
		index.NewController(catalogService, indexTemplate),
		catalog.NewController(catalogService, catalogTemplate),
		sell.NewController(sellService, sellTemplate),
		cart.NewController(cartService, cartTemplate),
	}
}
