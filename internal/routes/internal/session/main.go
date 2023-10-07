package session

import (
	"shorty-urls-server/internal/database"
	"time"

	"github.com/gofiber/fiber/v2/middleware/session"
)

var SessionStore *session.Store

func InitSession() {
	SessionStore = session.New(session.Config{
		Storage:    database.RedisClient,
		Expiration: time.Hour,
	})
}
