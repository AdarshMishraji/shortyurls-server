package app

import (
	"encoding/gob"
	"shorty-urls-server/internal/database"
	"shorty-urls-server/internal/internal/utils"
	"shorty-urls-server/internal/routes"

	"github.com/joho/godotenv"
)

func Start() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	gob.Register(utils.LocationInfo{})
	gob.Register(utils.DeviceInfo{})

	database.ConnectDB()
	database.ConnectRedis()
	routes.HandleRoutes()
}
