package middleware

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	// "mado/internal/config"
)

func corsMiddleware() gin.HandlerFunc {
	// cfg := config.Get()
	// cfg.CORS.AllowOrigins = append(cfg.CORS.AllowOrigins, "http://0.0.0.0:3000/")
	corsProvided := cors.New(cors.Config{
		// AllowAllOrigins: true,
		AllowOrigins: []string{"*"}, // Provide your list of allowed origins here
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowHeaders:     []string{"Origin", "Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
	})

	return corsProvided
}

// func corsMiddleware() gin.HandlerFunc {
// 	cfg := config.Get()

// 	corsProvided := cors.New(cors.Config{
// 		AllowOrigins: cfg.CORS.AllowOrigins,
// 		AllowMethods: []string{
// 			http.MethodGet,
// 			http.MethodPost,
// 			http.MethodPut,
// 			http.MethodDelete,
// 			http.MethodOptions,
// 		},
// 		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
// 		AllowCredentials: true,
// 	})

// 	return corsProvided
// }
