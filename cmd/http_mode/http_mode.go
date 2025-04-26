package http_mode

import (
	"custom-database/internal/http/handlers"
	"log"

	"github.com/gin-gonic/gin"
)

type QueryRequest struct {
	Query string `json:"query" binding:"required"`
}

type QueryResponse struct {
	Success bool        `json:"success"`
	Result  interface{} `json:"result,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func RunHttpServer(handlers handlers.HttpHandlers, port string) {
	router := gin.Default()

	router.Use(CorsMiddleware)

	router.POST("/query", handlers.HandleQuery)

	log.Printf("HTTP сервер запущен на порту %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Error starting HTTP server:", err)
	}
}
