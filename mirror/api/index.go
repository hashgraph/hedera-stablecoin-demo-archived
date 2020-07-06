package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.io/hashgraph/stable-coin/mirror/api/notification"
	"github.io/hashgraph/stable-coin/mirror/api/routes"
	"os"
)

func Run() {
	r := gin.New()

	r.Use(logger.SetLogger())
	r.Use(gin.Recovery())

	// https://github.com/gin-contrib/cors
	r.Use(cors.Default())
	r.GET("/v1/token", routes.GetToken)
	r.GET("/v1/token/userExists/:username", routes.GetUserExists)
	r.GET("/v1/token/balance/:address", routes.GetUserBalanceByAddress)
	r.GET("/v1/token/users/:address", routes.GetOtherUsersByAddress)
	r.GET("/v1/token/usersSearch/:username", routes.GetUsersByPartialMatch)
	r.GET("/v1/token/operations/:username", routes.GetUserOperationsByUsername)

	r.GET("/ws", notification.Handler)

	err := r.Run(":" + os.Getenv("MIRROR_PORT"))
	if err != nil {
		panic(err)
	}
}
