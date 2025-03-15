package main

import (
	"log"
	"net/http"

	"github.com/is-web-y26/m3302-milovatskiy/internal/deploy"
)

func main() {
	address := ":14041"

	log.Printf("server is running on %s", address)
	http.ListenAndServe(address, deploy.Setup())
}
