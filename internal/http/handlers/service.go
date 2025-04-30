package handlers

// "github.com/gin-gonic/gin"

type HttpHandlers interface {
	// HandleSqlQuery(c *gin.Context)
}

type handlers struct {
}

func NewHttpHandlers() HttpHandlers {
	return &handlers{}
}
