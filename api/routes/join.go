package routes

import (
	"github.com/labstack/echo/v4"
	"github.io/hashgraph/stable-coin/pb"
	"net/http"
)

func SendAnnounce(c echo.Context) error {
	var req struct {
		PublicKey string `json:"publicKey"`
		Username  string `json:"username"`
	}

	err := c.Bind(&req)
	if err != nil {
		return err
	}

	v := &pb.Join{
		Address:  req.PublicKey,
		Username: req.Username,
	}

	go sendTransaction(v, &pb.Primitive{Primitive: &pb.Primitive_Join{Join: v}})

	return c.JSON(http.StatusAccepted, transactionResponse{
		Status:  true,
		Message: "Join request sent",
	})
}
