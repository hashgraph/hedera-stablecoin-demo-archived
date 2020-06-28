package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"net/http/httputil"
	"net/url"
	"os"
)

var proxy *httputil.ReverseProxy

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	u, err := url.Parse(os.Getenv("MIRROR_ADDRESS"))
	if err != nil {
		panic(err)
	}

	proxy = httputil.NewSingleHostReverseProxy(u)
}

func GetUserExists(c *gin.Context) {
	proxy.ServeHTTP(c.Writer, c.Request)
}

func GetUserBalanceByAddress(c *gin.Context) {
	proxy.ServeHTTP(c.Writer, c.Request)
}
