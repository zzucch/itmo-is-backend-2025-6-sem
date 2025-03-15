package main

import (
	"net/http"

	"github.com/is-web-y26/m3302-milovatskiy/internal/deploy"
)

func main() {
	address := ":14041"

	http.ListenAndServe(address, deploy.Setup())
}
