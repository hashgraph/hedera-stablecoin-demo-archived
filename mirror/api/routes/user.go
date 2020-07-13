package routes

import (
	"encoding/hex"
	"github.com/hashgraph/hedera-sdk-go"
	"github.com/labstack/echo/v4"
	"github.io/hashgraph/stable-coin/mirror/state"
	"net/http"
	"strings"
)

func GetUsersByFrozenStatus(c echo.Context) error {
	type FrozenUser struct {
		Name  string   `json:"username"`
	}
	var frozen []FrozenUser

	isFrozen := c.Param("frozen")

	state.Frozen.Range(func(userNameI, frozenI interface{}) bool {
		userName := userNameI.(string)

		if (frozenI.(bool) == (isFrozen == "true")) && userName != "Admin" {
			frozen = append(frozen, FrozenUser{Name: userName})
		}

		return true
	})
	return c.JSON(http.StatusOK, frozen)
}

func GetUsersByPartialMatch(c echo.Context) error {
	userNames := []string{}

	searchValue := strings.ToUpper(c.Param("username"))

	if searchValue != "" {

		state.Address.Range(func(addressI, userNameI interface{}) bool {
			userName := userNameI.(string)

			if strings.Contains(strings.ToUpper(userName), searchValue) {
				userNames = append(userNames, userName)
			}

			return true
		})
	}
	if len(userNames) > 10 {
		return c.JSON(http.StatusOK, userNames[0:10])
	} else {
		return c.JSON(http.StatusOK, userNames)
	}
}

func IsAdminUser(c echo.Context) error {

	var adminUserI interface{}
	var exists bool

	valid := false

	adminHederaPublicKey, err := hedera.Ed25519PublicKeyFromString(c.Param("publicKey"))
	if err != nil {
		valid = false
	} else {
		adminPublicKeyHex := hex.EncodeToString(adminHederaPublicKey.Bytes())
		if adminUserI, exists = state.Address.Load(adminPublicKeyHex); exists {
			// key match, is it admin ?
			foundUserName := adminUserI.(string)
			if foundUserName == "Admin" {
				valid = true
			} else {
				valid = false
			}
		} else {
			// key not known
			valid = false
		}
	}

	return c.JSON(http.StatusOK, echo.Map{
		"valid": valid,
	})
}

//func GetOtherUsersByAddress(c echo.Context) error {
//	userNames := []string{}
//
//	// FIXME: This should be by username
//	hederaPublicKey, err := hedera.Ed25519PublicKeyFromString(c.Param("address"))
//	if err != nil {
//		return err
//	}
//
//	excludeAddress := hex.EncodeToString(hederaPublicKey.Bytes())
//
//	state.Address.Range(func(addressI, userNameI interface{}) bool {
//		address := addressI.(string)
//		userName := userNameI.(string)
//
//		if excludeAddress != address {
//			userNames = append(userNames, userName)
//		}
//
//		return true
//	})
//
//	return c.JSON(http.StatusOK, userNames)
//}
