package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetOtherUsersByAddress(c *gin.Context) {
	// address := c.Param("address")
	// TODO

	c.JSON(http.StatusOK, []interface{}{})
}
