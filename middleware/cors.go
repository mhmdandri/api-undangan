package middleware

import (
	"api-undangan/config"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	origins := []string{
		"http://localhost:3000",
		"http://localhost:5173",
		"http://127.0.0.1:3000",
		"http://127.0.0.1:5173",
		"http://192.168.1.10:3000",
		"https://weddingofandricica.me",
	}

	if config.Cfg != nil && config.Cfg.ClientOrigin != "" {
		originTrimmed := strings.TrimRight(config.Cfg.ClientOrigin, "/")
		found := false
		for _, o := range origins {
			if strings.EqualFold(o, originTrimmed) {
				found = true
				break
			}
		}
		if !found {
			origins = append(origins, originTrimmed)
		}
	}

	return cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
