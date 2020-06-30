package main

import (
	"github.com/gofiber/fiber"
	"github.com/joho/godotenv"
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
	app := fiber.New()

	//r.Use(logger.SetLogger())
	//r.Use(gin.Recovery())

	// https://github.com/gin-contrib/cors
	//r.Use(cors.Default())

	//r.GET("/v1/token", routes.GetToken)
	//r.GET("/v1/token/userExists/:username", routes.GetUserExists)
	//r.GET("/v1/token/balance/:address", routes.GetUserBalanceByAddress)
	//r.GET("/v1/token/users/:address", routes.GetOtherUsersByAddress)
	//r.GET("/v1/token/operations/:username", routes.GetUserOperationsByUsername)

	app.Post("/v1/token/join", routes.SendAnnounce)
	app.Post("/v1/token/mintTo", routes.SendMint)
	app.Post("/v1/token/transaction", routes.SendRawTransaction)

	//r.GET("/ws", notification.Handler)

	// TODO: Runs on :8080 by default but can be overridden by $PORT
	err := app.Listen(3128)
	if err != nil {
		panic(err)
	}
}
