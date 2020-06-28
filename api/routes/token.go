package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
)

var tokenDecimals int
var tokenSymbol string
var tokenName string

type tokenMeta struct {
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	Decimals int    `json:"decimals"`
}

func getTokenMeta() tokenMeta {
	if len(tokenName) == 0 {
		var err error
		tokenDecimals, err = strconv.Atoi(os.Getenv("DECIMALS"))
		if err != nil {
			panic(err)
		}

		tokenName = os.Getenv("TOKEN_NAME")
		tokenSymbol = os.Getenv("SYMBOL")
	}

	return tokenMeta{
		Name:     tokenName,
		Symbol:   tokenSymbol,
		Decimals: tokenDecimals,
	}
}

func GetToken(c *gin.Context) {
	c.JSON(http.StatusOK, getTokenMeta())
}
