package main

import (
	"log"
	"net/http"

	"github.com/is-web-y26/m3302-milovatskiy/internal/deploy"
)

func main() {
	address := ":14041"

	if err := deploy.Deploy(); err != nil {
		log.Fatal()
	}

	log.Printf("server is running on %s", address)
	if err := http.ListenAndServe(address, deploy.Setup()); err != nil {
		log.Fatal(err)
	}
}
