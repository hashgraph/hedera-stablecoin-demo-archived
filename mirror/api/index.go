package api

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/rs/zerolog"
	"github.com/ziflex/lecho/v2"
	"github.io/hashgraph/stable-coin/mirror/api/notification"
	"github.io/hashgraph/stable-coin/mirror/api/routes"
	"os"
)

func Run() {
	e := echo.New()

	logger := lecho.New(os.Stderr, lecho.WithTimestamp())
	logger.SetOutput(zerolog.ConsoleWriter{Out: os.Stderr, NoColor: false})

	// configure log level for mirror from env
	logger.SetLevel(mustParseEchoLogLevel(os.Getenv("MIRROR_API_LOG")))

	e.Logger = logger

	e.Use(middleware.Recover())

	e.Use(lecho.Middleware(lecho.Config{
		Logger: logger,
	}))

	e.Use(middleware.CORS())

	e.GET("/v1/token", routes.GetToken)
	e.GET("/v1/token/userExists/:username", routes.GetUserExists)
	e.GET("/v1/token/balance/:address", routes.GetUserBalanceByAddress)
	e.GET("/v1/token/users/:address", routes.GetOtherUsersByAddress)
	e.GET("/v1/token/operations/:username", routes.GetUserOperationsByUsername)
	e.GET("/v1/token/usersSearch/:username", routes.GetUsersByPartialMatch)

	e.GET("/ws", notification.Handler)

	err := e.Start(":" + os.Getenv("MIRROR_PORT"))
	if err != nil {
		panic(err)
	}
}

func mustParseEchoLogLevel(s string) log.Lvl {
	return map[string]log.Lvl{
		"DEBUG": log.DEBUG,
		"INFO":  log.INFO,
		"WARN":  log.WARN,
		"ERROR": log.ERROR,
		"OFF":   log.OFF,
	}[s]
}
