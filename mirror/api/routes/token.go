package routes

import (
	"github.com/labstack/echo/v4"
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

func GetToken(c echo.Context) error {
	return c.JSON(http.StatusOK, getTokenMeta())
}
