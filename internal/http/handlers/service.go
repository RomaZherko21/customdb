package handlers

import (
	"custom-database/internal/lexer"

	"github.com/gin-gonic/gin"
)

type HttpHandlers interface {
	HandleQuery(c *gin.Context)
}

type handlers struct {
	lexer lexer.Lexer
}

func NewHttpHandlers(lexer lexer.Lexer) HttpHandlers {
	return &handlers{
		lexer: lexer,
	}
}
