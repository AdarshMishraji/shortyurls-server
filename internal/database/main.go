package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gofiber/storage/redis/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

var RedisClient *redis.Storage

func ConnectDB() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")
	password := os.Getenv("DB_PASSWORD")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s TimeZone=Asia/Kolkata", host, user, password, dbName, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Enable debug mode
	})

	DB = db

	if err != nil {
		log.Fatal(err)
	}

	// Migrate()
	fmt.Println("Connected to database")
}

func Migrate() {
	DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")

	tables := []interface{}{
		new(User),
		new(UserLoginHistory),
		new(ShortenURL),
		new(ShortenURLVisit),
	}

	err := DB.AutoMigrate(tables...)
	if err != nil {
		log.Fatal(err)
	}
	err = DB.Exec("ALTER TABLE shorten_urls ADD CONSTRAINT chk_is_deleted_is_active CHECK ((is_deleted = FALSE AND (is_active = TRUE OR is_active = FALSE)) OR (is_deleted = TRUE AND is_active = FALSE));").Error
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			fmt.Println("Constraint already exists")
		} else {
			log.Fatal(err)
		}
	}
}

func ConnectRedis() {
	RedisClient = redis.New(redis.Config{
		URL: os.Getenv("REDIS_URL"),
	})

	_, err := RedisClient.Conn().Ping(context.Background()).Result()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to redis")
}
