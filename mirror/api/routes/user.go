package routes

import (
	"encoding/hex"
	"github.com/hashgraph/hedera-sdk-go"
	"github.com/labstack/echo/v4"
	"github.io/hashgraph/stable-coin/mirror/state"
	"net/http"
	"strings"
)

func GetUsersByPartialMatch(c *gin.Context) {
	userNames := []string{}

	searchValue := c.Param("username")

	if searchValue != "" {

		state.Address.Range(func(addressI, userNameI interface{}) bool {
			userName := userNameI.(string)

			if strings.Contains(userName, searchValue) {
				userNames = append(userNames, userName)
			}

			return true
		})
	}
	if len(userNames) > 10 {
		c.JSON(http.StatusOK, userNames[0:10])
	} else {
		c.JSON(http.StatusOK, userNames)
	}
}

func GetOtherUsersByAddress(c echo.Context) error {
	userNames := []string{}

	// FIXME: This should be by username
	hederaPublicKey, err := hedera.Ed25519PublicKeyFromString(c.Param("address"))
	if err != nil {
		return err
	}

	excludeAddress := hex.EncodeToString(hederaPublicKey.Bytes())

	state.Address.Range(func(addressI, userNameI interface{}) bool {
		address := addressI.(string)
		userName := userNameI.(string)

		if excludeAddress != address {
			userNames = append(userNames, userName)
		}

		return true
	})

	return c.JSON(http.StatusOK, userNames)
}
