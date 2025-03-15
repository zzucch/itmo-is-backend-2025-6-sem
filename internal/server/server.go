package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/is-web-y26/m3302-milovatskiy/internal/config"
)

func Start(logger *log.Logger, config *config.Config) error {
	address := fmt.Sprintf(":%s", config.Port)

	log.Printf("server is running on %s", address)
	http.ListenAndServe(address, Setup())

	return nil
}
