package handlers

import (
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

func (h *handlers) HandleQuery(c *gin.Context) {
	var request QueryRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	query := request.Query

	err := h.lexer.ParseQuery(query)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Query executed successfully"})
}
