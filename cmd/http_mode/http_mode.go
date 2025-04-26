package http_mode

import (
	"custom-database/internal/http/handlers"
	"log"

	_ "custom-database/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Custom Database API
// @version 1.0
// @description API для работы с кастомной базой данных
// @host localhost:8080
// @BasePath /
func RunHttpServer(handlers handlers.HttpHandlers, port string) {
	router := gin.Default()

	router.Use(CorsMiddleware)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.POST("/query", handlers.HandleSqlQuery)

	log.Printf("HTTP сервер запущен на порту %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Error starting HTTP server:", err)
	}
}
