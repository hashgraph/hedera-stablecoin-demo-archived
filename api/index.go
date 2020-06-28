package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.io/hashgraph/stable-coin/api/routes"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
}

func main() {
	r := gin.Default()

	// https://github.com/gin-contrib/cors
	r.Use(cors.Default())

	r.GET("/v1/token", routes.GetToken)

	// NOTE: Runs on :8080 by default but can be overridden by $PORT
	err := r.Run()
	if err != nil {
		panic(err)
	}
}
