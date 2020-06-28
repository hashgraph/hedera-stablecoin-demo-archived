package routes

import (
	"github.com/gin-gonic/gin"
	"github.io/hashgraph/stable-coin/data"
	"github.io/hashgraph/stable-coin/mirror/state"
	"net/http"
)

func GetUserExists(c *gin.Context) {
	var err error

	username := c.Param("username")
	var exists = true

	if _, ok := state.User[username]; !ok {
		exists, err = data.GetUserExists(username)
		if err != nil {
			panic(err)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"exists": exists,
	})
}
