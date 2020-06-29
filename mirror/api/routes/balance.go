package routes

import (
	"database/sql"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"github.com/hashgraph/hedera-sdk-go"
	"github.io/hashgraph/stable-coin/data"
	"github.io/hashgraph/stable-coin/mirror/state"
	"net/http"
)

func GetUserBalanceByAddress(c *gin.Context) {
	var err error

	address := c.Param("address")
	hederaPublicKey, err := hedera.Ed25519PublicKeyFromString(address)
	if err != nil {
		panic(err)
	}

	publicKeyHex := hex.EncodeToString(hederaPublicKey.Bytes())

	if username, ok := state.Address.Load(publicKeyHex); ok {
		if balance, ok := state.Balance.Load(username); ok {
			c.JSON(http.StatusOK, gin.H{
				"balance": balance,
			})

			return
		}
	}

	balance, _, err := data.GetUserBalanceByAddress(hederaPublicKey.Bytes())
	if err == sql.ErrNoRows {
		balance = 0
	} else if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"balance": balance,
	})
}
