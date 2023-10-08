package main

import (
	"shorty-urls-server/internal/app"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	app.Start()
}

// todo: handle error efficiently
