package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
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
	r := gin.New()

	r.Use(logger.SetLogger())
	r.Use(gin.Recovery())

	// https://github.com/gin-contrib/cors
	r.Use(cors.Default())

	r.GET("/v1/token/userExists/:username", routes.GetUserExists)
	r.GET("/v1/token/balance/:address", routes.GetUserBalanceByAddress)
	r.GET("/v1/token/users/:address", routes.GetOtherUsersByAddress)
	r.GET("/v1/token/operations/:username", routes.GetUserOperationsByUsername)

	r.POST("/v1/token/join", routes.SendAnnounce)
	r.POST("/v1/token/mintTo", routes.SendMint)
	r.POST("/v1/token/transaction", routes.SendRawTransaction)

	// NOTE: Runs on :8080 by default but can be overridden by $PORT
	err := r.Run()
	if err != nil {
		panic(err)
	}
}
