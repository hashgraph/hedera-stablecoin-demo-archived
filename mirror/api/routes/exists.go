package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
	"github.io/hashgraph/stable-coin/data"
	"github.io/hashgraph/stable-coin/mirror/state"
	"net/http"
)

func GetUserExists(c echo.Context) error {
	var err error

	username := c.Param("username")
	var exists = true

	if _, ok := state.User.Load(username); !ok {
		exists, err = data.GetUserExists(username)
		if err != nil {
			return err
		}
	}

	return c.JSON(http.StatusOK, gin.H{
		"exists": exists,
	})
}
