package main

import (
	"log"
	"os"
	"time"

	"github.com/is-web-y26/m3302-milovatskiy/internal/config"
	"github.com/is-web-y26/m3302-milovatskiy/internal/server"
)

func main() {
	config := config.GetDefaultConfig()
	logger := initLogger()

	if err := server.Start(logger, config); err != nil {
		log.Fatal(err)
	}
}

func initLogger() *log.Logger {
	const logDir = "./logs/"

	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	file, err := os.Create(logDir + time.Now().Format(time.DateTime) + ".log")
	if err != nil {
		log.Fatal(err)
	}

	const prefix = ""
	const flags = log.LstdFlags | log.LUTC | log.Lmicroseconds

	return log.New(file, prefix, flags)
}
