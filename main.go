package main

import (
	"shorty-urls-server/database"
	"shorty-urls-server/routes"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	database.ConnectDB()
	routes.HandleRoutes()
}

// todo: handle error efficiently
