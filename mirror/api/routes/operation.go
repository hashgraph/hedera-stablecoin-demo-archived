package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetUserOperationsByUsername(c *gin.Context) {
	// username := c.Param("username")
	// TODO

	c.JSON(http.StatusOK, gin.H{"operations": []interface{}{}})
}
