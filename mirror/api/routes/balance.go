package routes

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
	"github.io/hashgraph/stable-coin/data"
	"github.io/hashgraph/stable-coin/mirror/state"
	"net/http"
)

func GetUserBalanceByAddress(c echo.Context) error {
	var err error

	username := c.Param("username")
	if balance, ok := state.Balance.Load(username); ok {
		return c.JSON(http.StatusOK, gin.H{
			"balance": balance,
		})
	}
	balance, _, err := data.GetUserBalanceByUsername(username)

	if err == sql.ErrNoRows {
		balance = 0
	} else if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, gin.H{
		"balance": balance,
	})
}
