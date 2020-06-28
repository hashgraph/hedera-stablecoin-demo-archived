package routes

import (
	"github.com/gin-gonic/gin"
	"github.io/hashgraph/stable-coin/pb"
	"net/http"
)

func SendAnnounce(c *gin.Context) {
	var req struct {
		PublicKey string `json:"publicKey"`
		Username  string `json:"username"`
	}

	err := c.BindJSON(&req)
	if err != nil {
		return
	}

	v := &pb.Join{
		Address:  req.PublicKey,
		Username: req.Username,
	}

	sendTransaction(v, &pb.Primitive{Primitive: &pb.Primitive_Join{Join: v}})

	c.JSON(http.StatusAccepted, transactionResponse{
		Status:  true,
		Message: "Join request sent",
	})
}
