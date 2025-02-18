package main

import (
	"log"

	"github.com/is-web-y26/m3302-milovatskiy/internal/config"
	"github.com/is-web-y26/m3302-milovatskiy/internal/server"
)

func main() {
	config := config.GetDefaultConfig()

	if err := server.Start(config); err != nil {
		log.Fatal(err)
	}
}
