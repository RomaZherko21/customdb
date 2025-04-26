package mode

import (
	"custom-database/internal/lexer"
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

func RunHttpServer(lexer lexer.Lexer, port string) {
	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	router.POST("/query", func(c *gin.Context) {
		var request QueryRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, QueryResponse{
				Success: false,
				Error:   "Invalid request format: " + err.Error(),
			})
			return
		}

		err := lexer.ParseQuery(request.Query)
		if err != nil {
			c.JSON(400, QueryResponse{
				Success: false,
				Error:   err.Error(),
			})
			return
		}

		c.JSON(200, QueryResponse{
			Success: true,
			Result:  "Query executed successfully",
		})
	})

	log.Printf("HTTP сервер запущен на порту %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Error starting HTTP server:", err)
	}
}
