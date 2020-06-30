package main

import (
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.io/hashgraph/stable-coin/api/routes"
	"os"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, NoColor: false})
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// https://github.com/gin-contrib/cors
	e.Use(middleware.CORS())

	// TODO: e.GET("/v1/token", routes.GetToken)
	// TODO: e.GET("/v1/token/userExists/:username", routes.GetUserExists)
	// TODO: e.GET("/v1/token/balance/:address", routes.GetUserBalanceByAddress)
	// TODO: e.GET("/v1/token/users/:address", routes.GetOtherUsersByAddress)
	// TODO: e.GET("/v1/token/operations/:username", routes.GetUserOperationsByUsername)

	e.POST("/v1/token/join", routes.SendAnnounce)
	e.POST("/v1/token/mintTo", routes.SendMint)
	e.POST("/v1/token/transaction", routes.SendRawTransaction)

	// TODO: e.GET("/ws", notification.Handler)

	// NOTE: Runs on :8080 by default but can be overridden by $PORT
	err := e.Start(":" + os.Getenv("PORT"))
	if err != nil {
		panic(err)
	}
}
