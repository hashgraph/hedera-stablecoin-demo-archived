package routes

import (
	"database/sql"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"github.com/hashgraph/hedera-sdk-go"
	"github.com/labstack/echo/v4"
	"github.io/hashgraph/stable-coin/data"
	"github.io/hashgraph/stable-coin/mirror/state"
	"net/http"
)

func GetUserBalanceByAddress(c echo.Context) error {
	var err error

	address := c.Param("address")
	hederaPublicKey, err := hedera.Ed25519PublicKeyFromString(address)

	if err != nil {
		return err
	}

	publicKeyHex := hex.EncodeToString(hederaPublicKey.Bytes())

	if username, ok := state.Address.Load(publicKeyHex); ok {
		if balance, ok := state.Balance.Load(username); ok {
			return c.JSON(http.StatusOK, gin.H{
				"balance": balance,
			})
		}
	}

	balance, _, err := data.GetUserBalanceByAddress(hederaPublicKey.Bytes())

	if err == sql.ErrNoRows {
		balance = 0
	} else if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, gin.H{
		"balance": balance,
	})
}
