package routes

import (
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"github.com/hashgraph/hedera-sdk-go"
	"github.io/hashgraph/stable-coin/mirror/state"
	"net/http"
)

func GetOtherUsersByAddress(c *gin.Context) {
	userNames := make([]string, 0, len(state.User))

	// FIXME: This should be by username
	hederaPublicKey, err := hedera.Ed25519PublicKeyFromString(c.Param("address"))
	if err != nil {
		panic(err)
	}

	excludeAddress := hex.EncodeToString(hederaPublicKey.Bytes())

	for address, userName := range state.Address {
		if excludeAddress != address {
			userNames = append(userNames, userName)
		}
	}

	c.JSON(http.StatusOK, userNames)
}
