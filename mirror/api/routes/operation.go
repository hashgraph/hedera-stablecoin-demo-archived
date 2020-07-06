package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
	"github.io/hashgraph/stable-coin/data"
	"github.io/hashgraph/stable-coin/domain"
	"github.io/hashgraph/stable-coin/mirror/state"
	"net/http"
)

func GetUserOperationsByUsername(c echo.Context) error {
	username := c.Param("username")
	existingOperations, err := data.GetOperationsForUsername(username)
	if err != nil {
		return err
	}

	pendingOperations := state.GetPendingOperationsForUser(username)

	operations := make([]domain.Operation, 0, len(existingOperations)+len(pendingOperations))
	operations = append(operations, pendingOperations...)
	operations = append(operations, existingOperations...)

	return c.JSON(http.StatusOK, gin.H{"operations": operations})
}
