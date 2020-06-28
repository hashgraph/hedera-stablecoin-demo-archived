package api

import (
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.io/hashgraph/stable-coin/mirror/api/routes"
	"os"
)

func Run() {
	r := gin.New()

	r.Use(logger.SetLogger())
	r.Use(gin.Recovery())

	r.GET("/v1/token/userExists/:username", routes.GetUserExists)
	r.GET("/v1/token/balance/:address", routes.GetUserBalanceByAddress)

	err := r.Run(":" + os.Getenv("MIRROR_PORT"))
	if err != nil {
		panic(err)
	}
}
